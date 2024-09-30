package main

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sort"

	"log/slog"

	svg "github.com/ajstarks/svgo"
	"github.com/formicidae-tracker/tag-layouter/internal/tag"
	"github.com/jessevdk/go-flags"
	"golang.org/x/image/tiff"
)

// A LevelHandler wraps a Handler with an Enabled method
// that returns false for levels below a minimum.
type LevelHandler struct {
	level   slog.Leveler
	handler slog.Handler
}

// NewLevelHandler returns a LevelHandler with the given level.
// All methods except Enabled delegate to h.
func NewLevelHandler(level slog.Leveler, h slog.Handler) *LevelHandler {
	// Optimization: avoid chains of LevelHandlers.
	if lh, ok := h.(*LevelHandler); ok {
		h = lh.Handler()
	}
	return &LevelHandler{level, h}
}

// Enabled implements Handler.Enabled by reporting whether
// level is at least as large as h's level.
func (h *LevelHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle implements Handler.Handle.
func (h *LevelHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

// WithAttrs implements Handler.WithAttrs.
func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewLevelHandler(h.level, h.handler.WithAttrs(attrs))
}

// WithGroup implements Handler.WithGroup.
func (h *LevelHandler) WithGroup(name string) slog.Handler {
	return NewLevelHandler(h.level, h.handler.WithGroup(name))
}

// Handler returns the Handler wrapped by h.
func (h *LevelHandler) Handler() slog.Handler {
	return h.handler
}

type App struct {
	Blocks           []tag.FamilyBlock `short:"t" long:"family-and-size" description:"Families and size to use. format: 'name:size[:begin-end]'" required:"yes"`
	Columns          int               `long:"columns" description:"Number of column to display multiple families" default:"1"`
	TagBorder        float64           `long:"individual-tag-border" description:"border between tags in column layout" default:"0.2"`
	CutLineRatio     float64           `long:"cut-line-ratio" description:"ratio of the border between tags that should be a cut line" default:"0.0"`
	BlockMargin      float64           `long:"block-margin" description:"margin between families in mm" default:"2.0"`
	PaperSize        tag.Size          `short:"P" long:"paper-size" description:"Output paper size in mm. format: <width>x<height>" default:"210.0x297.0"`
	Margin           float64           `long:"margin" description:"Margnins of the paper in mm" default:"20.0"`
	LabelRoundedSize bool              `long:"label-rounded-size" description:"Label the rounded size instead of the actual size"`
	DPI              int               `short:"d" long:"dpi" description:"DPI to use" default:"2400"`
	Verbose          []bool            `short:"V" long:"verbose" description:"increase verbose level"`

	CPUProfile string `long:"cpuprofile" description:"profile CPU usage"`

	Args struct {
		File flags.Filename
	} `positional-args:"yes" required:"yes"`
}

func main() {
	if err := execute(); err != nil {
		log.Fatalf("unhandled error: %s", err)
	}
}

func execute() error {
	var app App
	if _, err := flags.Parse(&app); err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		}
		os.Exit(1)
	}
	closeCPU, err := app.SetUpProfile()
	if err != nil {
		return err
	}
	defer closeCPU()

	app.SetupLogger()

	columns, _, err := app.Layout()
	if err != nil {
		return err
	}

	drawer, closeImage, err := app.CreateDrawer()

	for _, column := range columns {
		for _, block := range column.Blocks {
			label := block.LabelWithDesiredSize()
			if app.LabelRoundedSize == true {
				label = block.LabelWithSize(tag.PixelToMM(app.DPI, block.ActualTagWidth))
			}
			block.Render(drawer, label)
		}
	}

	return closeImage()
}

func (a App) SetupLogger() {
	var level slog.Level
	switch len(a.Verbose) {
	case 0:
		level = slog.LevelWarn
	case 1:
		level = slog.LevelInfo
	default:
		level = slog.LevelDebug
	}
	defaultHandler := slog.Default().Handler()
	slog.SetDefault(slog.New(NewLevelHandler(level, defaultHandler)))
	log.SetOutput(os.Stdout)
}

func (a App) Layout() ([]Column, int64, error) {
	if a.Columns < 1 {
		return nil, 0, fmt.Errorf("invalid number of column %d: must be >=1", a.Columns)
	}

	blockMarginDot := tag.MMToPixel(a.DPI, a.BlockMargin)
	marginDot := tag.MMToPixel(a.DPI, a.Margin)
	slog.Debug("paper margins", "pixel", marginDot, "mm", a.Margin)
	slog.Debug("block margins", "pixel", blockMarginDot, "mm", a.BlockMargin)

	columnWidthDot := tag.MMToPixel(a.DPI, a.PaperSize.Width) - 2*marginDot
	columnWidthDot -= blockMarginDot * (a.Columns - 1)
	columnWidthDot /= a.Columns

	columnHeightDot := tag.MMToPixel(a.DPI, a.PaperSize.Height) - 2*marginDot

	fullWidth := []PlacedBlock{}
	incompleteWidth := []PlacedBlock{}
	for _, f := range a.Blocks {
		pf := a.ComputeFamilySize(f, columnWidthDot)
		slog.Debug("block inner dimension in pixels",
			"block", f.String(),
			"tag", pf.ActualTagWidth,
			"border", pf.ActualBorderWidth,
			"cutline", pf.CutLineWidth,
		)
		if pf.Width < columnWidthDot {
			incompleteWidth = append(incompleteWidth, pf)
		} else {
			fullWidth = append(fullWidth, pf)
		}
	}

	sort.Slice(incompleteWidth, func(i, j int) bool {
		return incompleteWidth[i].Width < incompleteWidth[j].Width
	})
	sort.Slice(fullWidth, func(i, j int) bool {
		return fullWidth[i].Height > fullWidth[j].Height
	})

	columns := make([]Column, a.Columns)
	for i := range columns {
		columns[i].XOffset = i*(columnWidthDot+blockMarginDot) + marginDot
		columns[i].Width = 0
		columns[i].Height = 0
		columns[i].LastRowHeight = 0
	}

	for _, block := range fullWidth {
		fitted := false
		for idxCol, _ := range columns {
			if (block.Height + columns[idxCol].Height + blockMarginDot) > columnHeightDot {
				continue
			}
			if len(columns[idxCol].Blocks) == 0 {
				columns[idxCol].Height = marginDot - blockMarginDot
			}
			slog.Info("block placement",
				"block", block.String(),
				"column", idxCol,
				"row", len(columns[idxCol].Blocks),
			)
			block.X = columns[idxCol].XOffset
			block.Y = columns[idxCol].Height + blockMarginDot

			columns[idxCol].Blocks = append(columns[idxCol].Blocks, block)
			columns[idxCol].Height += block.Height + blockMarginDot
			fitted = true
			break
		}

		if fitted == false {
			return nil, 0, fmt.Errorf("Could not fit %s in layout", block)
		}
	}

	for _, block := range incompleteWidth {
		fitted := false
		for idxCol, _ := range columns {
			if (block.Height + columns[idxCol].Height + blockMarginDot) > columnHeightDot {
				//not fitting in height anyway
				continue
			}
			//if we are building a new line
			//check if it fits on the same line
			if (block.Width + columns[idxCol].Width) > columnWidthDot {
				//no so we terminate the line
				columns[idxCol].Width = 0
				columns[idxCol].Height = columns[idxCol].LastRowHeight
				columns[idxCol].LastRowHeight = 0
				//we recheck if we can be put in height
				if block.Height+columns[idxCol].Height+blockMarginDot > columnHeightDot {
					continue
				}
			}

			slog.Info("placing small block",
				"block", block.String(),
				"column", idxCol,
				"row", len(columns[idxCol].Blocks),
			)

			block.X = columns[idxCol].XOffset + columns[idxCol].Width
			block.Y = columns[idxCol].Height + blockMarginDot

			columns[idxCol].Blocks = append(columns[idxCol].Blocks, block)
			columns[idxCol].Width += block.Width + blockMarginDot
			columns[idxCol].LastRowHeight = max(columns[idxCol].LastRowHeight, columns[idxCol].Height+blockMarginDot+block.Height)
			fitted = true
			break
		}
		if fitted == false {
			return nil, 0, fmt.Errorf("Could not fill %s in layout", block)
		}
	}
	var N int64 = 0
	for i, col := range columns {
		for j, block := range col.Blocks {
			N += int64(block.FamilyBlock.Len())
			trueTagSizeMM := tag.PixelToMM(a.DPI, block.ActualTagWidth)
			error := math.Abs(trueTagSizeMM-block.SizeMM) / block.SizeMM * 100.0
			fmt.Printf("block %s placed in column %d at position %d\n", block, i, j)
			fmt.Printf("---> tag range:  %s, %d tags\n", block.FamilyBlock, block.FamilyBlock.Len())
			fmt.Printf("---> tag size: %dpx %ddot.px⁻¹ %.3gmm error: %.2f%%\n",
				block.Family.TotalWidth,
				block.ActualTagWidth/block.Family.TotalWidth,
				trueTagSizeMM,
				error,
			)
			fmt.Printf("---> tag in between space: %d px %.3g mm\n",
				block.ActualBorderWidth,
				tag.PixelToMM(a.DPI, block.ActualBorderWidth),
			)
			fmt.Printf("---> cutline size: %d px %.3g mm\n",
				block.CutLineWidth,
				tag.PixelToMM(a.DPI, block.CutLineWidth),
			)
		}
	}

	return columns, N, nil
}

func (a App) ComputeFamilySize(block tag.FamilyBlock, columnWidthDot int) PlacedBlock {
	actualTagWidth, actualBorderWidth, cutlineWidth := PerfectPixelSizeMM(a.DPI,
		block.SizeMM, a.TagBorder, a.CutLineRatio, block.Family.TotalWidth)

	skips := len(block.LabelWithDesiredSize())*1/2 + 1
	nbSlots := block.Len() + skips

	nbTagsPerRow := (columnWidthDot - actualBorderWidth) / (actualTagWidth + actualBorderWidth)
	nbRows := nbSlots / nbTagsPerRow
	if nbSlots%nbTagsPerRow != 0 {
		nbRows += 1
	}

	height := nbRows*(actualTagWidth+actualBorderWidth) + actualBorderWidth
	width := 0
	if nbRows > 1 {
		width = columnWidthDot
	} else {
		width = nbSlots*(actualTagWidth+actualBorderWidth) + actualBorderWidth
	}

	return PlacedBlock{
		FamilyBlock:       block,
		Height:            height,
		Width:             width,
		X:                 0,
		Y:                 0,
		ActualTagWidth:    actualTagWidth,
		ActualBorderWidth: actualBorderWidth,
		CutLineWidth:      cutlineWidth,
		NTagsPerRow:       nbTagsPerRow,
		Skips:             skips,
		DPI:               a.DPI,
	}
}

func (app App) CreateDrawer() (Drawer, func() error, error) {
	ext := filepath.Ext(string(app.Args.File))
	switch ext {
	case ".jpg":
		return app.createJPEGDrawer()
	case ".jpeg":
		return app.createJPEGDrawer()
	case ".tiff":
		return app.createTIFFDrawer()
	case ".png":
		return app.createPNGDrawer()
	case ".svg":
		return app.createSVGDrawer()
	default:
		return nil, nil, fmt.Errorf("unsupported extension '%s' for filepath '%s'", ext, app.Args.File)
	}

}

func (app App) createImageDrawer(encoder func(w io.Writer, img image.Image) error) (Drawer, func() error, error) {
	file, err := os.Create(string(app.Args.File))
	if err != nil {
		return nil, nil, err
	}

	drawer, err := NewImageDrawer(app.PaperSize.Width, app.PaperSize.Height, app.DPI)
	if err != nil {
		return nil, nil, err
	}

	encodeAndClose := func() error {
		var errs []error
		errs = append(errs, encoder(file, drawer.(*imageDrawer).img))
		errs = append(errs, file.Close())
		return errors.Join(errs...)
	}
	return drawer, encodeAndClose, nil
}

func (app App) createJPEGDrawer() (Drawer, func() error, error) {
	return app.createImageDrawer(func(w io.Writer, img image.Image) error {
		return jpeg.Encode(w, img, nil)
	})
}

func (app App) createTIFFDrawer() (Drawer, func() error, error) {
	return app.createImageDrawer(func(w io.Writer, img image.Image) error {
		return tiff.Encode(w, img, nil)
	})
}

func (app App) createPNGDrawer() (Drawer, func() error, error) {
	return app.createImageDrawer(func(w io.Writer, img image.Image) error {
		return png.Encode(w, img)
	})
}

func (app App) createSVGDrawer() (VectorDrawer, func() error, error) {
	file, err := os.Create(string(app.Args.File))
	if err != nil {
		return nil, nil, err
	}
	svg := svg.New(file)
	drawer := NewSVGDrawer(svg, app.PaperSize.Width, app.PaperSize.Height, app.DPI, false)
	close := func() error {
		svg.End()
		return file.Close()
	}
	return drawer, close, nil
}

func (app App) SetUpProfile() (func() error, error) {
	if len(app.CPUProfile) == 0 {
		return func() error { return nil }, nil
	}
	f, err := os.Create(app.CPUProfile)
	if err != nil {
		return nil, fmt.Errorf("could not create cpuprofile file '%s': %w", app.CPUProfile, err)
	}
	pprof.StartCPUProfile(f)
	return func() error {
		pprof.StopCPUProfile()
		return nil
	}, nil
}
