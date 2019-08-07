package main

import (
	"fmt"
	"image/color"
	"io"
	"math"
	"os"

	svg "github.com/ajstarks/svgo"
)

type toDotConverter func(float64) int

type SVGDrawer struct {
	Dotter
	f   io.Closer
	SVG *svg.SVG
}

func NewSVGDrawer(filepath string, width, height float64, DPI int) (Drawer, error) {
	f, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	res := &SVGDrawer{
		Dotter: Dotter{float64(DPI)},
		f:      f,
		SVG:    svg.New(f),
	}

	res.SVG.Start(res.ToDot(width), res.ToDot(height))

	return res, nil
}

func (d *SVGDrawer) Close() error {
	d.SVG.End()
	return d.f.Close()
}

func (d *SVGDrawer) RotateTranslate(x, y int, angle float64) {
	radians := angle * math.Pi / 180.0
	xCanvas := float64(x)*math.Cos(radians) + float64(y)*math.Sin(radians)
	yCanvas := -float64(x)*math.Sin(radians) + float64(y)*math.Cos(radians)

	d.SVG.RotateTranslate(int(xCanvas), int(yCanvas), angle)
}

func (d *SVGDrawer) EndRotateTranslate() {
	d.SVG.Gend()
}

func (d *SVGDrawer) DrawRectangle(x, y, w, h int, c color.Color) {
	r, g, b, _ := c.RGBA()
	d.SVG.Rect(x, y, w, h, fmt.Sprintf("stroke:none;fill:rgb(%d,%d,%d)", r/256, g/256, b/256))
}

func (d *SVGDrawer) DrawLine(x1, y1, x2, y2, border int, c color.Color) {
	r, g, b, _ := c.RGBA()
	d.SVG.Line(x1, y1, x2, y2, fmt.Sprintf("stroke:rgb(%d,%d,%d);stroke-width:%d", r/256, g/256, b/256, border))
}

func (d *SVGDrawer) Label(x, y int, height int, label string, c color.Color) float64 {
	d.SVG.Text(x, y+height, label, fmt.Sprintf("font-size:%d;font-family:Roboto", height))
	return 0

}

func (d *SVGDrawer) DrawCircle(x, y, radius, border int, c color.Color) {
	r, g, b, _ := c.RGBA()
	d.SVG.Circle(x, y, radius, fmt.Sprintf("stroke:rgb(%d,%d,%d);fill:none;stroke-width:%d", r/256, g/256, b/256, border))
}
