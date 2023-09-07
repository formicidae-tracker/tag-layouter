package main

import (
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"math"

	"github.com/formicidae-tracker/tag-layouter/internal/tag"
	"github.com/schollz/progressbar/v3"
)

type ColumnLayouter struct {
	Width            float64
	Height           float64
	NColumns         int
	FamilyMargin     float64
	TagBorder        float64
	PaperBorder      float64
	CutLine          float64
	LabelroundedSize bool
	DPI              int
}

type PlacedBlock struct {
	tag.FamilyBlock
	Height int
	Width  int
	X      int
	Y      int

	ActualTagWidth    int
	ActualBorderWidth int
	CutLineWidth      int
	NTagsPerRow       int
	Skips             int
	DPI               int
}

func PerfectPixelSizeMM(DPI int, sizeMM float64, borderRatio float64, cutlineRatio float64, pixelSize int) (tagSizeDot int, borderSizeDot int, cutLineSizeDot int) {
	perfectPixelSize := tag.MMToPixel(DPI, sizeMM/float64(pixelSize))
	//check if a little bit bigger is not better
	errors := []float64{
		math.Abs(tag.PixelToMM(DPI, perfectPixelSize*pixelSize) - sizeMM),
		math.Abs(tag.PixelToMM(DPI, (perfectPixelSize+1)*pixelSize) - sizeMM),
	}
	if errors[1] < errors[0] {
		perfectPixelSize += 1
	}

	tagSizeDot = perfectPixelSize * pixelSize
	borderSizeDot = int(math.Round(float64(tagSizeDot) * borderRatio))

	cutLineSizeDot = int(math.Round(float64(borderSizeDot) * cutlineRatio))
	if cutlineRatio != 0.0 {
		if cutLineSizeDot == 0 {
			cutLineSizeDot = 1
		}
		if (borderSizeDot-cutLineSizeDot)%2 == 1 {
			borderSizeDot += 1
		}
	}

	return tagSizeDot, borderSizeDot, cutLineSizeDot
}

type Column struct {
	Blocks        []PlacedBlock
	XOffset       int
	Width         int
	Height        int
	LastRowHeight int
}

func (b PlacedBlock) Render(drawer Drawer, label string) error {

	pb := progressbar.Default(int64(b.FamilyBlock.Len()),
		fmt.Sprintf("rendering %s", b))

	scale := b.ActualTagWidth / b.Family.TotalWidth
	if b.ActualTagWidth%b.Family.TotalWidth != 0 {
		return fmt.Errorf("invalid scaling ratio %d / %d: should be divisible",
			b.ActualTagWidth, b.Family.TotalWidth)
	}

	ix := b.Skips % b.NTagsPerRow
	iy := b.Skips / b.NTagsPerRow

	cutLinePos := (b.ActualBorderWidth - b.CutLineWidth) / 2
	isFirst := true
	for _, r := range b.Ranges {
		end := r.End
		if end < 0 {
			end = len(b.Family.Codes)
		}
		for i := r.Begin; i < end; i++ {
			b.renderTag(drawer, i, ix, iy, scale, cutLinePos, isFirst)
			pb.Add(1)
		}
	}

	pb.Close()
	drawer.Label(image.Pt(b.X+b.ActualBorderWidth, b.Y+b.ActualBorderWidth/2),
		label, b.ActualTagWidth, color.Gray{0})

	return nil
}

func (b PlacedBlock) renderTag(drawer Drawer, idx, ix, iy, scale, cutLinePos int, isFirst bool) {
	vectorDrawer, ok := drawer.(VectorDrawer)

	slog.Debug("rendering tag", "index", idx, "code", b.Family.Codes[idx])
	var pos image.Point
	pos.X = ix*(b.ActualTagWidth+b.ActualBorderWidth) + b.X + b.ActualBorderWidth
	pos.Y = iy*(b.ActualTagWidth+b.ActualBorderWidth) + b.Y + b.ActualBorderWidth

	img := b.Family.RenderTag(idx)

	drawer.TranslateScale(pos, scale)
	if ok == true {
		vectorDrawTag(vectorDrawer, img)
	} else {
		drawTag(drawer, img)
	}
	drawer.EndTranslate()

	ix += 1
	if ix >= b.NTagsPerRow {
		ix = 0
		iy += 1
	}
	if b.CutLineWidth == 0 || true {
		return
	}

	drawer.DrawRectangle(
		image.Rect(pos.X+b.ActualTagWidth+cutLinePos, pos.Y,
			b.CutLineWidth, b.ActualTagWidth),
		color.Gray{127})

	drawer.DrawRectangle(
		image.Rect(pos.X, pos.Y+b.ActualTagWidth+cutLinePos,
			b.ActualTagWidth, b.CutLineWidth),
		color.Gray{127})

	if ix == 1 || isFirst == true {
		isFirst = false
		drawer.DrawRectangle(
			image.Rect(pos.X-b.CutLineWidth-cutLinePos, pos.Y,
				b.CutLineWidth, b.ActualTagWidth),
			color.Gray{127})
	}

	if iy == 0 || iy == 1 && ix <= b.Skips {
		drawer.DrawRectangle(
			image.Rect(pos.X, pos.Y-cutLinePos-b.CutLineWidth,
				b.ActualTagWidth, b.CutLineWidth),
			color.Gray{127})
	}

}
