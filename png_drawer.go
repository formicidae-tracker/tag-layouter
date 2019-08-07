package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/tiff"
)

type imageEncoder func(io.Writer, *image.Gray) error

var encoders map[string]imageEncoder

func matchEncoder(filename string) (imageEncoder, error) {
	ext := filepath.Ext(filename)
	if res, ok := encoders[ext]; ok == true {
		return res, nil
	}
	return nil, fmt.Errorf("Unsupported file extension '%s'", ext)
}

func init() {
	encoders = make(map[string]imageEncoder)

	encoders[".tiff"] = func(w io.Writer, i *image.Gray) error {
		return tiff.Encode(w, i, &tiff.Options{
			Compression: tiff.Deflate,
			Predictor:   true,
		})
	}
	encoders[".tif"] = encoders[".tiff"]

	encoders[".png"] = func(w io.Writer, i *image.Gray) error {
		return png.Encode(w, i)
	}

}

type ImageDrawer struct {
	Dotter
	f    *os.File
	data *image.Gray

	encoder imageEncoder
	fd      *freetype.Context
	x       []int
	y       []int
}

func NewImageDrawer(filepath string, width, height float64, DPI int) (Drawer, error) {
	f, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	encoder, err := matchEncoder(filepath)
	if err != nil {
		return nil, err
	}

	res := &ImageDrawer{
		Dotter:  Dotter{float64(DPI)},
		f:       f,
		encoder: encoder,
	}
	res.data = image.NewGray(image.Rect(0, 0, res.ToDot(width), res.ToDot(height)))

	monofont, err := truetype.Parse(gomono.TTF)
	if err != nil {
		return nil, err
	}
	res.fd = freetype.NewContext()
	res.fd.SetDPI(float64(DPI))
	res.fd.SetFont(monofont)
	res.fd.SetClip(res.data.Bounds())
	res.fd.SetDst(res.data)

	return res, nil
}

func (d *ImageDrawer) Close() error {
	err := d.encoder(d.f, d.data)
	if err != nil {
		return err
	}
	return d.f.Close()
}

func (d *ImageDrawer) Offsets() (int, int) {
	if len(d.x) == 0 {
		return 0, 0
	}
	return d.x[len(d.x)-1], d.y[len(d.y)-1]
}

type Vec2 struct {
	x, y float64
}

func (v *Vec2) Norm() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (d *ImageDrawer) DrawCircle(x, y, r, hb int, c color.Color) {
	xo, yo := d.Offsets()
	for i := x - r - hb - 1; i <= x+r+hb+1; i++ {
		for j := y - r - hb - 1; j <= y+r+hb+1; j++ {
			dist := (&Vec2{float64(i - x), float64(j - y)}).Norm()
			if dist < float64(r+hb) && dist > float64(r-hb) {
				d.data.Set(i+xo, j+yo, c)
			}
		}
	}

}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (d *ImageDrawer) DrawLine(x1, y1, x2, y2, b int, c color.Color) {
	xo, yo := d.Offsets()
	if abs(x1-x2) > (y1 - y2) {
		if x1 > x2 {
			x1, x2 = x2, x1
			y1, y2 = y2, y1
		}

		for x := x1; x <= x2; x++ {
			y := y1 + (y2-y1)*(x-x1)/(x2-x1)
			d.data.Set(x+xo, y+yo, c)
		}
	} else {
		if y1 > y2 {
			x1, x2 = x2, x1
			y1, y2 = y2, y1
		}

		for y := y1; y <= y2; y++ {
			x := x1 + (x2-x1)*(y-y1)/(y2-y1)
			d.data.Set(x+xo, y+yo, c)
		}

	}
}

func (d *ImageDrawer) DrawRectangle(x, y, w, h int, c color.Color) {
	xo, yo := d.Offsets()
	g, _ := color.GrayModel.Convert(c).(color.Gray)
	pv := 0
	for i := x + xo; i < x+xo+w; i++ {
		pv += 1
		if pv%1000 == 0 {
			log.Printf("Large rectangle col %d/%d", i-xo, w)
		}
		for j := y + yo; j < y+yo+h; j++ {
			d.data.SetGray(i, j, g)
		}
	}
}

func (d *ImageDrawer) RotateTranslate(x, y int, angle float64) {
	if angle != 0 {
		panic("angles are not supported")
	}
	d.x = append(d.x, x)
	d.y = append(d.y, y)
}

func (d *ImageDrawer) EndRotateTranslate() {
	if len(d.x) == 0 {
		return
	}
	d.x = d.x[0:(len(d.x) - 1)]
	d.y = d.y[0:(len(d.y) - 1)]
}

const (
	PtInMM float64 = 0.352778
)

func (d *ImageDrawer) Label(x, y int, height int, label string, c color.Color) float64 {
	xo, yo := d.Offsets()

	d.fd.SetSrc(&image.Uniform{c})
	size := d.ToMM(height) / PtInMM
	d.fd.SetFontSize(size)
	pt := freetype.Pt(x+xo, y+yo+int(d.fd.PointToFixed(size)>>6))

	advance, _ := d.fd.DrawString(label, pt)

	// Binarizing rasterized font
	for i := pt.X.Ceil() + xo; i < advance.X.Ceil()+xo; i++ {
		for j := y + yo; j <= pt.Y.Ceil()+1+yo; j++ {
			if d.data.GrayAt(i, j).Y < 200 {
				d.data.Set(i, j, color.Black)
			} else {
				d.data.Set(i, j, color.White)
			}
		}
	}

	return d.ToMM((advance.X.Ceil() - x))

}
