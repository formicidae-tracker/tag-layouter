package main

import (
	"fmt"
	"image/color"
	"math/rand"
)

type ArenaLayouter struct {
	Border float64
	Number int
	Width  float64
	Height float64
}

func (l *ArenaLayouter) Layout(drawer Drawer, families []FamilyBlock) error {
	if len(families) != 1 {
		return fmt.Errorf("Arena layouter only supports a single family (got:%d)", len(families))
	}
	set := map[int]Point{}

	if l.Border < 0 {
		return fmt.Errorf("Border cannot be negative")
	}

	if l.Border > 0.0 {
		drawer.DrawRectangle(drawer.ToDot(l.Border/2), drawer.ToDot(l.Border/2), drawer.ToDot(l.Width-l.Border), drawer.ToDot(l.Height-l.Border), color.Gray{Y: 200})
		drawer.DrawRectangle(drawer.ToDot(l.Border), drawer.ToDot(l.Border), drawer.ToDot(l.Width-2*l.Border), drawer.ToDot(l.Height-2*l.Border), color.White)
	}

	for i := 0; i < l.Number; i++ {
		angle := rand.Float64() * 360.0
		idx := 0
		for {
			idx = rand.Intn(len(families[0].Family.Codes) - 1)
			if _, ok := set[idx]; ok == true {
				continue
			}
			break
		}
		x := 0.0
		y := 0.0

		for {
			x = rand.Float64()*(l.Width-2*l.Border-2*families[0].Size) + l.Border + families[0].Size
			y = rand.Float64()*(l.Height-2*l.Border-2*families[0].Size) + l.Border + families[0].Size
			p := Point{x, y}
			if Touches(set, p, families[0].Size*3) == true {
				continue
			}
			set[idx] = p
			break
		}

		DrawTag(drawer, families[0].Family, families[0].Family.Codes[i], x, y, families[0].Size, angle, &i)
	}
	return nil
}
