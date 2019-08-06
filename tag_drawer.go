package main

import (
	"fmt"
	"image/color"
	"math"
)

func DrawTag(drawer Drawer, blackBorder, whiteBorder, codeSize int, payload int64, x, y, size, angle float64, value *int) error {

	sizeInPX := codeSize + 2*blackBorder + 2*whiteBorder
	pixelSize := drawer.ToDot(size / float64(sizeInPX))

	if pixelSize == 0.0 {
		return fmt.Errorf("tag size too small")
	}

	drawer.RotateTranslate(drawer.ToDot(x), drawer.ToDot(y), angle)
	defer drawer.EndRotateTranslate()

	if value != nil {
		sizeInDot := pixelSize * sizeInPX
		lInDot := 3 * sizeInDot / 2
		hInDot := int(math.Sqrt(3) / 2.0 * float64(lInDot))
		x1 := sizeInDot / 2
		y1 := sizeInDot/2 - 2*hInDot/3
		x2 := sizeInDot/2 - lInDot/2
		y2 := sizeInDot/2 + hInDot/3
		x3 := sizeInDot/2 + lInDot/2
		y3 := y2

		drawer.DrawLine(x1, y1, x2, y2, 3, color.Black)
		drawer.DrawLine(x1, y1, x3, y3, 3, color.Black)
		drawer.DrawLine(x2, y2, x3, y3, 3, color.Black)
		drawer.Label(sizeInDot/2+hInDot/2, sizeInDot/2, size/3, fmt.Sprintf("%d", *value), color.Black)

	}

	drawer.DrawRectangle(0, 0, sizeInPX*pixelSize, sizeInPX*pixelSize, color.White)

	drawer.DrawRectangle(pixelSize*whiteBorder, pixelSize*whiteBorder, (sizeInPX-2*whiteBorder)*pixelSize, (sizeInPX-2*whiteBorder)*pixelSize, color.Black)

	offset := (whiteBorder + blackBorder) * pixelSize
	for i := uint(0); i < uint(codeSize*codeSize); i++ {
		j := int(i) % codeSize
		k := (int(i) - j) / codeSize
		if (1<<i)&payload == 0 {
			continue
		}

		drawer.DrawRectangle(offset+(codeSize-1-j)*pixelSize, offset+(codeSize-1-k)*pixelSize, pixelSize, pixelSize, color.White)
	}

	return nil
}
