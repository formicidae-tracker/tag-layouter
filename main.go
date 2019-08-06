package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/jessevdk/go-flags"
)

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
	File   string  `short:"f" long:"file" description:"File to output" required:"true"`
	Size   float64 `short:"s" long:"size" description:"Tag Size" default:"1.6"`
	Family string  `long:"family" description:"Family to use" default:"36h11"`
	Number int     `short:"n" long:"number" description:"Number of tags" default:"20"`
	Width  float64 `short:"W" long:"width" description:"Width to use" default:"210"`
	Height float64 `short:"H" long:"height" description:"Height to use" default:"297"`
	Border float64 `short:"b" long:"border" description:"Draw a border" default:"20.0"`
	DPI    int     `short:"d" long:"dpi" description:"DPI to use" default:"2400"`
}

func Execute() error {
	opts := Options{}
	if _, err := flags.Parse(&opts); err != nil {
		return err
	}

	drawer, err := NewSVGDrawer(opts.File, opts.Width, opts.Height, opts.DPI)
	if err != nil {
		return err
	}
	defer drawer.Close()

	family, err := GetFamily(opts.Family)

	if err != nil {
		return err
	}

	set := map[int]Point{}

	if opts.Border < 0 {
		return fmt.Errorf("Border cannot be negative")
	}

	if opts.Border > 0.0 {
		drawer.DrawRectangle(drawer.ToDot(opts.Border/2), drawer.ToDot(opts.Border/2), drawer.ToDot(opts.Width-opts.Border), drawer.ToDot(opts.Height-opts.Border), color.Gray{Y: 200})
		drawer.DrawRectangle(drawer.ToDot(opts.Border), drawer.ToDot(opts.Border), drawer.ToDot(opts.Width-2*opts.Border), drawer.ToDot(opts.Height-2*opts.Border), color.White)
	}

	for i := 0; i < opts.Number; i++ {
		angle := rand.Float64() * 360.0
		idx := 0
		for {
			idx = rand.Intn(len(family.Codes) - 1)
			if _, ok := set[idx]; ok == true {
				continue
			}
			break
		}
		x := 0.0
		y := 0.0

		for {
			x = rand.Float64()*(opts.Width-2*opts.Border-2*opts.Size) + opts.Border + opts.Size
			y = rand.Float64()*(opts.Height-2*opts.Border-2*opts.Size) + opts.Border + opts.Size
			p := Point{x, y}
			if Touches(set, p, opts.Size*3) == true {
				continue
			}
			set[idx] = p
			break
		}

		DrawTag(drawer, family, family.Codes[i], x, y, opts.Size, angle, &i)
	}

	return nil

}

func main() {
	if err := Execute(); err != nil {
		log.Printf("Unhandled error: %s", err)
		os.Exit(1)
	}

}
