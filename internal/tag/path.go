package tag

import (
	"image"
)

var (
	up    = image.Point{0, -1}
	left  = image.Point{-1, 0}
	down  = image.Point{0, 1}
	right = image.Point{1, 0}
)

var relativeDirections = map[image.Point]struct {
	Left, Right image.Point
}{
	up:    {left, right},
	right: {up, down},
	down:  {right, left},
	left:  {down, up},
}

var vertexOffset = map[image.Point]image.Point{
	up:    image.Point{0, 0},
	right: image.Point{1, 0},
	down:  image.Point{1, 1},
	left:  image.Point{0, 1},
}

func BuildPath(img *image.Gray) [][]image.Point {

	vertices := []image.Point{}
	// 1. Find first foreground pixel
	background := uint8(0xff)
	pos := func() image.Point {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
				if img.GrayAt(x, y).Y != background {
					return image.Pt(x, y)
				}
			}
		}
		return img.Rect.Max
	}()

	direction := up
	pointIsForeground := func(p image.Point) bool {
		return p.In(img.Rect.Bounds()) && img.GrayAt(p.X, p.Y).Y != background
	}

	for {
		rel := relativeDirections[direction]

		front := pos.Add(direction)
		leftFront := front.Add(rel.Left)
		frontIsForeground := pointIsForeground(front)
		leftFrontIsForeground := pointIsForeground(leftFront)
		if leftFrontIsForeground == false && frontIsForeground == true {
			// we advance along the direction
			pos = front
			continue
		}

		newVertex := pos.Add(vertexOffset[direction])
		if len(vertices) > 0 && vertices[0] == newVertex {
			break
		}
		vertices = append(vertices, newVertex)

		if leftFrontIsForeground == true {
			//turn left
			direction = rel.Left
			pos = leftFront
		} else {
			//turn right
			direction = rel.Right
			//we stay at the same position.
		}

	}

	return [][]image.Point{vertices}
}
