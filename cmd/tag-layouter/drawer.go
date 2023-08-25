package main

import (
	"image"
	"image/color"

	"gihtub.com/formicidae-tracker/tag-layouter/internal/tag"
	svg "github.com/ajstarks/svgo"
)

type Drawer interface {
	TranslateScale(image.Point, int)
	EndTranslate()

	DrawRectangle(image.Rectangle, color.Color)
	Label(p image.Point, s string, size int, c color.Color)
}

type VectorDrawer interface {
	DrawPolygons([]tag.Polygon)
}

func drawTag(drawer Drawer, img *image.Gray) {
	for i, v := range img.Pix {
		if v == 0xff {
			continue
		}
		y := i / img.Stride
		x := i - y
		drawer.DrawRectangle(image.Rect(x, y, 1, 1), color.Black)
	}
}

func vectorDrawTag(drawer VectorDrawer, img *image.Gray) {
	polygons := tag.BuildPolygons(img)
	drawer.DrawPolygons(polygons)
}

type svgDrawer *svg.SVG

type imageDrawer struct {
	img       *image.Gray
	scales    []int
	positions []image.Point
}

func (d *imageDrawer) TranslateScale(pos image.Point, scale int) {
	d.scales = append(d.scales, scale)
	d.positions = append(d.positions, pos)
}

func (d *imageDrawer) EndTranslate(pos image.Point, scale int) {
	if min(len(d.scales), len(d.positions)) == 0 {
		return
	}
	d.scales = d.scales[:len(d.scales)-1]
	d.positions = d.positions[:len(d.scales)]
}

func (d *imageDrawer) DrawRectange(r image.Rectangle, c color.Color) {
	var pos image.Point
	scale := 1
	if min(len(d.scales), len(d.positions)) > 1 {
		pos = d.positions[len(d.positions)-1]
		scale = d.scales[len(d.scales)-1]
	}

	r = r.Add(pos)
	for y := scale * r.Min.Y; y < scale*r.Max.Y; y++ {
		for x := scale * r.Min.X; x < scale*r.Max.X; x++ {

		}
	}

}
