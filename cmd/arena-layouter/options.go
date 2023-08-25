package main

import (
	"fmt"
	"math"
	"math/rand"

	"gihtub.com/formicidae-tracker/tag-layouter/internal/tag"
	svg "github.com/ajstarks/svgo"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Family    tag.FamilyBlock `short:"t" long:"family-and-size" description:"Family and size to use. format: 'name:size:begin-end'" required:"yes"`
	Number    int             `short:"n" long:"number" description:"Number of tags to display in an arena" required:"yes" `
	PaperSize tag.Size        `short:"P" long:"paper-size" description:"Output size to use in mm" default:"210.0x297.0"`
	ArenaSize tag.Size        `short:"A" long:"arena-size" description:"Arena size to use in mm" default:"180x220"`
	DPI       int             `short:"d" long:"dpi" description:"DPI to use" default:"300"`

	Args struct {
		File flags.Filename
	} `positional-args:"yes" required:"yes"`
}

func (o Options) LayoutArena(SVG *svg.SVG) error {

	if math.Trunc(o.PaperSize.Height) != o.PaperSize.Height ||
		math.Trunc(o.PaperSize.Width) != o.PaperSize.Width {
		return fmt.Errorf("invalid paper size '%s': sub-millimeter paper size are not supported", o.PaperSize)
	}

	SVG.StartviewUnit(int(o.PaperSize.Width), int(o.PaperSize.Height), "mm",
		0, 0, int(o.PaperSize.Width), int(o.PaperSize.Height))
	defer SVG.End()

	fmt.Fprintln(SVG.Writer, `<rect width="100%" height="100%" style="fill:#fff" />`)

	if err := o.layoutBorder(SVG); err != nil {
		return err
	}

	if err := o.layoutTags(SVG); err != nil {
		return err
	}

	return nil

}

func (o Options) layoutBorder(SVG *svg.SVG) error {

	if o.ArenaSize.Width >= o.PaperSize.Width ||
		o.ArenaSize.Height >= o.PaperSize.Height {
		return fmt.Errorf("incompatible paper size (%s) and arena size (%s): the arena must fit on the paper",
			o.PaperSize, o.ArenaSize)
	}

	xMin := (o.PaperSize.Width - o.ArenaSize.Width) / 2
	xMax := xMin + o.ArenaSize.Width
	yMin := (o.PaperSize.Height - o.ArenaSize.Height) / 2
	yMax := yMin + o.ArenaSize.Height

	border := min(xMin, yMin) / 2
	if border <= 0.0 {
		return nil
	}

	points := []tag.PointF[float64]{
		{X: xMin - border, Y: yMin - border},
		{X: xMax + border, Y: yMin - border},
		{X: xMax + border, Y: yMax + border},
		{X: xMin - border, Y: yMax + border},
		{X: xMin, Y: yMin},
		{X: xMax, Y: yMin},
		{X: xMax, Y: yMax},
		{X: xMin, Y: yMax},
	}

	SVG.Path(tag.BuildSVGD(points[:4])+" "+tag.BuildSVGD(points[4:]), `style="fill:#7f7f7f"`)

	return nil
}

func (o Options) layoutTags(SVG *svg.SVG) error {
	if o.Number > len(o.Family.Family.Codes) {
		return fmt.Errorf("invalid number of tag %d: it must be smaller than the number of tag in '%s' (%d)",
			o.Number, o.Family.Family.Name, len(o.Family.Family.Codes))
	}

	indexes := rand.Perm(len(o.Family.Family.Codes))[:o.Number]

	positions := make([]tag.PointF[float64], 0, o.Number)

	scale := o.Family.SizeMM / tag.PixelToMM(o.DPI, o.Family.Family.TotalWidth)

	touchesAny := func(p tag.PointF[float64]) bool {
		radius2 := 3 * o.Family.SizeMM
		radius2 *= radius2
		for _, pos := range positions {
			dX := pos.X - p.X
			dY := pos.Y - p.Y
			dist2 := dX*dX + dY*dY
			if dist2 < radius2 {
				return true
			}
		}
		return false
	}

	for i := 0; i < o.Number; i++ {
		var pos tag.PointF[float64]
		for {
			pos.X = rand.Float64() * (o.ArenaSize.Width - o.Family.SizeMM)
			pos.Y = rand.Float64() * (o.ArenaSize.Height - o.Family.SizeMM)
			if touchesAny(pos) == false {
				break
			}
		}
		positions = append(positions, pos)

		o.layoutTag(SVG, indexes[i], pos, scale)
	}

	return nil
}

func (o Options) layoutTag(SVG *svg.SVG, idx int, pos tag.PointF[float64], scale float64) {
	angle := rand.Float64() * 360.0

	SVG.Gtransform(fmt.Sprintf("translate(%g,%g),rotate(%g),scale(%f)", pos.X, pos.Y, angle, scale))
	defer SVG.Gend()

	tag.RenderToSVG(SVG, tag.BuildPolygons(o.Family.Family.RenderTag(idx)))
}
