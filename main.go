package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
)

//go:generate make -C apriltag
//go:generate make -C oldtags

type Point struct {
	X, Y float64
}

func Touches(points map[int]Point, toTest Point, radius float64) bool {
	for _, p := range points {
		distX := p.X - toTest.X
		distY := p.Y - toTest.Y
		dist := math.Sqrt(distX*distX + distY*distY)
		if dist < radius {
			return true
		}
	}
	return false
}

type Options struct {
	File             string   `short:"f" long:"file" description:"File to output" required:"true"`
	FamilyAndSize    []string `short:"t" long:"family-and-size" description:"Families and size to use. format: 'name:size:begin-end'"`
	ColumnNumber     int      `long:"column-number" description:"Number of column to display multiple families" default:"0"`
	TagBorder        float64  `long:"individual-tag-border" description:"border between tags in column layout" default:"0.2"`
	CutLineRatio     float64  `long:"cut-line-ratio" description:"ratio of the border between tags that should be a cut line" default:"0.0"`
	FamilyMargin     float64  `long:"family-margin" description:"margin between families in mm" default:"2.0"`
	ArenaNumber      int      `long:"arena-number" description:"Number of tags to display in an arena" default:"0"`
	Width            float64  `short:"W" long:"width" description:"Width to use" default:"210"`
	Height           float64  `short:"H" long:"height" description:"Height to use" default:"297"`
	PaperBorder      float64  `long:"paper-border" description:"Border for arena or paper" default:"20.0"`
	LabelRoundedSize bool     `long:"label-rounded-size" description:"Label the rounded size instead of the actual size"`
	DPI              int      `short:"d" long:"dpi" description:"DPI to use" default:"2400"`
}

func ExtractFamilyAndSizes(list []string) ([]FamilyBlock, error) {
	res := []FamilyBlock{}
	for _, fAndSize := range list {
		fargs := strings.Split(fAndSize, ":")
		if len(fargs) <= 1 {
			return res, fmt.Errorf("invalid family specification '%s': need at list family and size in the form '<name>:<size>'", fAndSize)
		}
		if len(fargs) > 3 {
			return res, fmt.Errorf("invalid family specification '%s':  expected '<name>:<size>:<range>'", fAndSize)
		}
		tf, err := GetFamily(fargs[0])
		if err != nil {
			return res, err
		}
		s, err := strconv.ParseFloat(fargs[1], 64)
		if err != nil {
			return res, err
		}

		if len(fargs) == 2 {
			res = append(res, FamilyBlock{
				Family: tf,
				Size:   s,
				Ranges: []Range{
					Range{
						Begin: 0,
						End:   len(tf.Codes),
					},
				},
			})
			continue
		}

		ranges, err := ExtractRanges(fargs[2])
		if err != nil {
			return res, err
		}
		if len(ranges) == 0 {
			return res, fmt.Errorf("Range for '%s' cannot be empty", fAndSize)
		}
		for i, r := range ranges {
			if r.Begin >= len(tf.Codes) {
				return res, fmt.Errorf("%d is out of range for %s in '%s'", r.Begin, fargs[0], fargs[2])
			}
			if r.End < 0 {
				ranges[i].End = len(tf.Codes)
			}
			if ranges[i].End > len(tf.Codes) {
				return res, fmt.Errorf("%d is out of range for %s in '%s'", ranges[i].End, fargs[0], fargs[2])
			}
		}

		res = append(res, FamilyBlock{
			Family: tf,
			Size:   s,
			Ranges: ranges,
		})
	}
	return res, nil
}

func Execute() error {
	opts := Options{}
	if _, err := flags.Parse(&opts); err != nil {
		return err
	}

	var drawer Drawer = nil
	var err error
	if filepath.Ext(opts.File) == ".svg" {
		drawer, err = NewSVGDrawer(opts.File, opts.Width, opts.Height, opts.DPI)
	} else {
		drawer, err = NewImageDrawer(opts.File, opts.Width, opts.Height, opts.DPI)
	}
	if err != nil {
		return err
	}
	defer drawer.Close()

	families, err := ExtractFamilyAndSizes(opts.FamilyAndSize)
	if err != nil {
		return err
	}

	var layouter Layouter = nil

	if opts.ArenaNumber != 0 && opts.ColumnNumber == 0 {
		layouter = &ArenaLayouter{
			Border: opts.PaperBorder,
			Number: opts.ArenaNumber,
			Width:  opts.Width,
			Height: opts.Height,
		}
	} else if opts.ColumnNumber != 0 && opts.ArenaNumber == 0 {
		layouter = &ColumnLayouter{
			Width:            opts.Width,
			Height:           opts.Height,
			NColumns:         opts.ColumnNumber,
			PaperBorder:      opts.PaperBorder,
			FamilyMargin:     opts.FamilyMargin,
			TagBorder:        opts.TagBorder,
			LabelroundedSize: opts.LabelRoundedSize,
			CutLine:          opts.CutLineRatio,
		}
	} else if opts.ColumnNumber != 0 && opts.ArenaNumber != 0 {
		return fmt.Errorf("Please specify either a column or either an arena layout")
	}
	if layouter == nil {
		return fmt.Errorf("Please specify a layout with either --arena-number or -- col-number")
	}

	err = layouter.Layout(drawer, families)
	if err != nil {
		log.Fatalf("Cannot layout : %s", err)
	}
	return nil
}

func main() {
	if err := Execute(); err != nil {
		log.Printf("Unhandled error: %s", err)
		os.Exit(1)
	}

}
