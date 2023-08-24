package main

import (
	"fmt"
	"math"
	"path/filepath"

	"gihtub.com/formicidae-tracker/tag-layouter/internal/tag"
	svg "github.com/ajstarks/svgo"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Family    tag.FamilyBlock `short:"t" long:"family-and-size" description:"Family and size to use. format: 'name:size:begin-end'" required:"yes"`
	Number    int             `short:"n" long:"number" description:"Number of tags to display in an arena" required:"yes" `
	PaperSize tag.Size        `short:"P" long:"paper-size" description:"Output size to use in mm" default:"210.0x297.0"`
	ArenaSize tag.Size        `short:"A" long:"arena-size" description:"Arena size to use in mm" default:"180x220"`
	Margin    float64         `short:"m" long:"margin" description:"For paper in mm " default:"20.0"`
	DPI       int             `short:"d" long:"dpi" description:"DPI to use" default:"300"`

	Args struct {
		File flags.Filename
	} `positional-args:"yes" required:"yes"`
}

func (o Options) Validate() error {
	if filepath.Ext(string(o.Args.File)) != ".svg" {
		return fmt.Errorf("invalid filepath '%s': only SVG are supported, filepath must end with '.svg'", o.Args.File)
	}

	if o.ArenaSize.Width >= o.PaperSize.Width ||
		o.ArenaSize.Height >= o.PaperSize.Height {
		return fmt.Errorf("incompatible paper size (%s) and arena size (%s): the arena must fit on the paper",
			o.PaperSize, o.ArenaSize)
	}

	if math.Trunc(o.PaperSize.Height) != o.PaperSize.Height ||
		math.Trunc(o.PaperSize.Width) != o.PaperSize.Width {
		return fmt.Errorf("invalid paper size '%s': sub-millimeter paper size are not supported", o.PaperSize)
	}

	return nil
}

func (o Options) LayoutArena(SVG *svg.SVG) {
	SVG.Startunit(int(o.PaperSize.Width), int(o.PaperSize.Height), "mm",
		`style="background-color:#fff;"`)
	defer SVG.End()

}
