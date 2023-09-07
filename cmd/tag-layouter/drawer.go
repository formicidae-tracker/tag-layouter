package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"log/slog"

	svg "github.com/ajstarks/svgo"
	"github.com/formicidae-tracker/tag-layouter/internal/tag"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/image/font/gofont/gomono"
)

//go:generate mockery --name Drawer

type Drawer interface {
	TranslateScale(image.Point, int)
	EndTranslate()

	DrawRectangle(image.Rectangle, color.Gray)
	Label(p image.Point, s string, size int, c color.Gray)
}

//go:generate mockery --name VectorDrawer

type VectorDrawer interface {
	Drawer
	DrawPolygons([]tag.Polygon)
}

func drawTag(drawer Drawer, img *image.Gray) {
	for i, v := range img.Pix {
		if v == 0xff {
			continue
		}
		y := i / img.Stride
		x := i - y*img.Stride
		drawer.DrawRectangle(image.Rect(x, y, x+1, y+1), color.Gray{})
	}
}

func vectorDrawTag(drawer VectorDrawer, img *image.Gray) {
	polygons := tag.BuildPolygons(img)
	drawer.DrawPolygons(polygons)
}

type imageDrawer struct {
	img       *image.Gray
	scales    []int
	positions []image.Point

	DPI int
	fd  *freetype.Context
}

func NewImageDrawer(width, height float64, DPI int) (Drawer, error) {
	slog.Info("new image drawer", "width", width, "height", height, "DPI", DPI)
	img := image.NewGray(image.Rect(0, 0,
		tag.MMToPixel(DPI, width), tag.MMToPixel(DPI, height)))

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	pb := progressbar.Default(100, "filling background")
	for i := 0; i < 100; i++ {
		rect := image.Rect(0, int(float64(h)*(float64(i)/100.0)),
			w, int(float64(h)*(float64(i+1)/100)))
		draw.Draw(img, rect, image.NewUniform(color.White), image.Point{}, draw.Src)

		pb.Add(1)
	}
	pb.Close()

	monofont, err := truetype.Parse(gomono.TTF)
	if err != nil {
		return nil, err
	}
	fd := freetype.NewContext()
	fd.SetDPI(float64(DPI))
	fd.SetFont(monofont)
	fd.SetClip(img.Bounds())
	fd.SetDst(img)

	return &imageDrawer{
		img: img,
		DPI: DPI,
		fd:  fd,
	}, nil
}

func (d *imageDrawer) TranslateScale(pos image.Point, scale int) {
	d.scales = append(d.scales, scale)
	d.positions = append(d.positions, pos)
	slog.Debug("image translate scale", "position", pos, "scale", scale,
		"positions", d.positions, "scales", d.scales)
}

func (d *imageDrawer) EndTranslate() {
	if min(len(d.scales), len(d.positions)) == 0 {
		return
	}
	d.scales = d.scales[:len(d.scales)-1]
	d.positions = d.positions[:len(d.scales)]
}

func (d *imageDrawer) DrawRectangle(r image.Rectangle, c color.Gray) {
	slog.Debug("drawing", "rectangle", r, "color", c)
	for i := min(len(d.scales), len(d.positions)) - 1; i >= 0; i-- {
		// apply stack of translate + scale (we have to apply in reverse order)
		pos := d.positions[i]
		scale := d.scales[i]
		slog.Debug("applying", "position", pos, "scale", scale)

		r.Min.X *= scale
		r.Min.Y *= scale
		r.Max.X *= scale
		r.Max.Y *= scale
		r = r.Add(pos)
	}
	slog.Debug("after transform", "rectangle", r)

	draw.Draw(d.img, r, image.NewUniform(c), image.Point{}, draw.Src)

}

func (d *imageDrawer) Label(p image.Point, s string, size int, c color.Gray) {
	slog.Debug("labelling", "label", s, "size", size, "at", p, "color", c)
	d.fd.SetSrc(image.NewUniform(c))
	fontHeightPt := 72 / 25.4 * tag.PixelToMM(d.DPI, size)
	d.fd.SetFontSize(fontHeightPt)

	pt := freetype.Pt(p.X, p.Y+int(d.fd.PointToFixed(fontHeightPt)>>6))
	advance, _ := d.fd.DrawString(s, pt)

	for i := pt.X.Ceil(); i < advance.X.Ceil(); i++ {
		for j := p.Y; j <= pt.Y.Ceil()+1; j++ {
			if d.img.GrayAt(i, j).Y < 200 {
				d.img.Set(i, j, color.Black)
			} else {
				d.img.Set(i, j, color.White)
			}
		}
	}
}

type svgDrawer svg.SVG

func (d *svgDrawer) TranslateScale(p image.Point, scale int) {
	(*svg.SVG)(d).Gtransform(fmt.Sprintf("translate(%d,%d),scale(%d)", p.X, p.Y, scale))
}

func (d *svgDrawer) EndTranslate() {
	(*svg.SVG)(d).Gend()
}

func (d *svgDrawer) DrawRectangle(r image.Rectangle, c color.Gray) {
	(*svg.SVG)(d).Rect(r.Min.X, r.Min.Y, r.Dx(), r.Dy(),
		fmt.Sprintf(`style="fill:%s"`, tag.ColorToHex(c)))
}

func (d *svgDrawer) Label(p image.Point, s string, size int, c color.Gray) {
	(*svg.SVG)(d).Text(p.X, p.Y+size, s, fmt.Sprintf("font-size:%dpx;font-family:Roboto Mono", size))
}

func (d *svgDrawer) DrawPolygons(polygons []tag.Polygon) {
	tag.RenderToSVG((*svg.SVG)(d), polygons)
}

func NewSVGDrawer(SVG *svg.SVG, width, height float64, DPI int, useMM bool) VectorDrawer {
	if useMM == true {
		SVG.StartviewUnit(int(width), int(height), "mm", 0, 0, int(width), int(height))
		SVG.Rect(0, 0, int(width), int(height), `style="fill:white"`)
	} else {
		w := tag.MMToPixel(DPI, width)
		h := tag.MMToPixel(DPI, height)
		SVG.Startview(w, h,
			0, 0, w, h)
		SVG.Rect(0, 0, w, h, `style="fill:white"`)
	}

	return (*svgDrawer)(SVG)
}
