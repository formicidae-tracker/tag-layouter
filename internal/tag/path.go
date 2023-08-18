package tag

import (
	"image"
	"image/color"
)

type direction int

const (
	up    direction = 0
	right direction = 1
	down  direction = 2
	left  direction = 3
)

var relativeVectors = []struct {
	Left, Front, Right image.Point
}{
	{image.Pt(-1, 0), image.Pt(0, -1), image.Pt(1, 0)}, // up
	{image.Pt(0, -1), image.Pt(1, 0), image.Pt(0, 1)},  // right
	{image.Pt(1, 0), image.Pt(0, 1), image.Pt(-1, 0)},  // down
	{image.Pt(0, 1), image.Pt(-1, 0), image.Pt(0, -1)}, // left
}

func turnLeft(d direction) direction {
	return (d - 1 + 4) % 4
}

func turnRight(d direction) direction {
	return (d + 1) % 4
}

var vertexOffset = []image.Point{
	image.Pt(0, 0), //up
	image.Pt(1, 0), //right
	image.Pt(1, 1), // down
	image.Pt(0, 1), //left
}

type pathBuilder struct {
	image, visited *image.Gray
}

type _color int

const (
	black _color = 0
	white _color = 1
)

func (b *pathBuilder) At(p image.Point) _color {
	if p.In(b.image.Rect) == false {
		return white
	}

	if b.image.GrayAt(p.X, p.Y).Y == 0x00 {
		return black
	}
	return white
}

func (b *pathBuilder) Visited(p image.Point) bool {
	if p.In(b.visited.Rect) == false {
		return true
	}

	if b.visited.GrayAt(p.X, p.Y).Y == 0x00 {
		return false
	}

	return true
}

func (b *pathBuilder) MarkVisited(p image.Point) {
	if p.In(b.visited.Rect) == false {
		return
	}
	b.visited.SetGray(p.X, p.Y, color.Gray{Y: 1})
}

func (b *pathBuilder) FindFirstUpEdge(start image.Point) image.Point {
	for y := max(b.image.Rect.Min.Y, start.Y); y < b.image.Rect.Max.Y; y++ {
		for x := max(b.image.Rect.Min.X, start.X); x < b.image.Rect.Max.X; x++ {
			p := image.Pt(x, y)
			if b.Visited(p) == true {
				continue
			}

			prev := p.Add(image.Pt(-1, 0))
			b.MarkVisited(p)
			if b.At(prev) != b.At(p) {
				return p
			}
		}
	}
	return b.image.Rect.Max
}

func (b *pathBuilder) buildPath(pos image.Point) []image.Point {
	insideColor := b.At(pos)
	direction := up
	res := []image.Point{}
	for {
		rel := relativeVectors[direction]
		front := pos.Add(rel.Front)
		frontIsInside := b.At(front) == insideColor
		frontLeft := front.Add(rel.Left)
		frontLeftIsInside := b.At(frontLeft) == insideColor
		b.MarkVisited(pos)
		b.MarkVisited(pos.Add(rel.Left))

		if !frontLeftIsInside && frontIsInside {
			//we advance along the direction
			pos = front
			continue
		}
		newVertex := pos.Add(vertexOffset[direction])
		if len(res) > 0 && res[0] == newVertex {
			return res
		}
		res = append(res, newVertex)
		if frontLeftIsInside && frontIsInside {
			direction = turnLeft(direction)
			pos = frontLeft
			b.MarkVisited(front) // we need to mark it because of the jump
		} else {
			direction = turnRight(direction)
		}
	}

}

func (b *pathBuilder) buildPaths() [][]image.Point {
	res := [][]image.Point{}

	pos := b.FindFirstUpEdge(b.image.Rect.Min)

	for pos != b.image.Rect.Max {
		res = append(res, b.buildPath(pos))
		pos = b.FindFirstUpEdge(pos)
	}

	return res
}

func newPathBuilder(img *image.Gray) *pathBuilder {
	return &pathBuilder{
		image:   img,
		visited: image.NewGray(img.Rect),
	}
}

func BuildPath(img *image.Gray) [][]image.Point {
	return newPathBuilder(img).buildPaths()
}
