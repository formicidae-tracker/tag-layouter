package tag

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"
)

// A Polygon is defined by its foreground color and the ordered list
// of its vertices.
type Polygon struct {
	Color    color.Color
	Vertices []image.Point
}

// Implements fmt.Stringer
func (p Polygon) String() string {
	return fmt.Sprintf("{ Color: %s, Vertices: %s}", p.Color, p.Vertices)
}

// BuildPolygons takes a monochromatic image, and returns the polygon
// that compose it. Its assumes background is white and foreground is
// black, as for a printer that would deposite black ink on white
// paper. Thus any pixel outside of the list of polygons should be
// considered to be rendered white.
func BuildPolygons(img *image.Gray) []Polygon {
	return newPathBuilder(img).buildPaths()
}

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
	b.visited.Pix[b.visited.PixOffset(p.X, p.Y)] += 1
}

func (b *pathBuilder) FindFirstUpEdge(start image.Point) image.Point {
	for y := max(b.image.Rect.Min.Y, start.Y); y < b.image.Rect.Max.Y; y++ {
		for x := max(b.image.Rect.Min.X, start.X); x < b.image.Rect.Max.X; x++ {
			start.X = b.image.Rect.Min.X
			//pos := image.Point{x, y}
			//log.Printf("Looking at %s", pos)
			//b.formatVisited(log.Writer(), pos)
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

type drawState struct {
	Black   bool
	Current bool
}

var colorCode = map[drawState]string{
	{Black: false, Current: false}: "30;47",
	{Black: false, Current: true}:  "30;46",
	{Black: true, Current: false}:  "37;40",
	{Black: true, Current: true}:   "36;40",
}

func (b *pathBuilder) formatVisited(w io.Writer, pos image.Point) {
	fmt.Fprintln(w, "┌"+strings.Repeat("─", 2*b.visited.Rect.Dx())+"┐")
	for y := 0; y < b.visited.Rect.Dy(); y++ {
		line := "│"
		for x := 0; x < b.visited.Rect.Dx(); x++ {

			v := b.visited.GrayAt(x, y).Y
			isCurrent := pos == image.Point{x, y}
			isBlack := b.image.GrayAt(x, y).Y == 0xff

			line += fmt.Sprintf("\033[%sm%2d", colorCode[drawState{Black: isBlack, Current: isCurrent}], v)

		}
		line += "\033[m│"
		fmt.Fprintln(w, line)
	}
	fmt.Fprintln(w, "└"+strings.Repeat("─", 2*b.visited.Rect.Dx())+"┘")

}

func (b *pathBuilder) buildPath(pos image.Point) []image.Point {
	insideColor := b.At(pos)
	direction := up
	res := []image.Point{}
	for {
		//log.Printf("%s", res)
		//b.formatVisited(log.Writer(), pos)
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

func (b *pathBuilder) buildPaths() []Polygon {
	res := []Polygon{}

	pos := b.FindFirstUpEdge(b.image.Rect.Min)

	for pos != b.image.Rect.Max {
		c := color.Black
		if b.image.GrayAt(pos.X, pos.Y).Y != 0x00 {
			c = color.White
		}
		res = append(res, Polygon{Color: c, Vertices: b.buildPath(pos)})
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
