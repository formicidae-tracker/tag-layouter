package main

import "image/color"

type Drawer interface {
	DrawRectangle(x, y, w, h int, c color.Color)
	RotateTranslate(x, y int, r float64)
	EndRotateTranslate()
	DrawLine(x1, y1, x2, y2, b int, c color.Color)
	Label(x, y int, heightInMM float64, label string, c color.Color) float64
	DrawCircle(x, y, r, b int, c color.Color)
	Close() error
	ToDot(float64) int
}
