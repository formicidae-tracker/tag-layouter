package main

import (
	"fmt"
	"image/color"
	"math"
)

func drawTagDotPrivate(drawer Drawer, tf *TagFamily, payload uint64, size int) error {
	pixelSize := size / tf.TotalWidth
	colorOut := color.White
	colorIn := color.Black
	if tf.ReversedBorder == true {
		colorOut, colorIn = colorIn, colorOut
	}

	drawer.DrawRectangle(0, 0, tf.TotalWidth*pixelSize, tf.TotalWidth*pixelSize, colorOut)
	offset := (tf.TotalWidth - tf.WidthAtBorder) / 2
	drawer.DrawRectangle(pixelSize*offset, pixelSize*offset, tf.WidthAtBorder*pixelSize, tf.WidthAtBorder*pixelSize, colorIn)

	for i := uint64(0); i < uint64(tf.NBits); i++ {
		bit := (uint64(1) << (uint64(tf.NBits) - 1 - i))
		isSet := payload&bit != 0
		backgroundIsBlack := tf.Inside[i] != tf.ReversedBorder

		if isSet != backgroundIsBlack {
			continue
		}

		colorPixel := color.White
		if backgroundIsBlack == false {
			colorPixel = color.Black
		}

		drawer.DrawRectangle((offset+tf.LocationX[i])*pixelSize, (offset+tf.LocationY[i])*pixelSize, pixelSize, pixelSize, colorPixel)
	}
	return nil
}

func DrawTagDot(drawer Drawer, tf *TagFamily, payload uint64, x, y, size int) error {
	drawer.RotateTranslate(x, y, 0.0)
	defer drawer.EndRotateTranslate()
	return drawTagDotPrivate(drawer, tf, payload, size)
}

func DrawTag(drawer Drawer, tf *TagFamily, payload uint64, x, y, size, angle float64, value *int) error {

	sizeInPX := tf.TotalWidth
	pixelSize := drawer.ToDot(size / float64(sizeInPX))

	if pixelSize == 0.0 {
		return fmt.Errorf("tag size too small")
	}

	drawer.RotateTranslate(drawer.ToDot(x), drawer.ToDot(y), angle)
	defer drawer.EndRotateTranslate()

	if value != nil {
		sizeInDot := pixelSize * sizeInPX
		lInDot := 3 * sizeInDot
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
		drawer.Label(sizeInDot/2+hInDot/2, sizeInDot/2, sizeInDot/3, fmt.Sprintf("%d", *value), color.Black)
	}

	return drawTagDotPrivate(drawer, tf, payload, tf.TotalWidth*pixelSize)
}
