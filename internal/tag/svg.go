package tag

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	svg "github.com/ajstarks/svgo"
	"golang.org/x/exp/constraints"
)

func toHex(c color.Color) string {
	r, g, b, a := c.RGBA()
	if a == 0xffff {
		return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", r>>8, g>>8, b>>8, a>>8)
}

func RenderToSVG(SVG *svg.SVG, polygons []Polygon) {
	if len(polygons) == 0 {
		return
	}
	paths := make([]string, len(polygons))

	for i := range polygons {
		paths[i] = BuildSVGPathData(polygons[i].Vertices)
	}

	SVG.Path(strings.Join(paths, " "),
		fmt.Sprintf("style=\"fill:%s;fill-rule:evenodd\"", toHex(polygons[0].Color)))
}

type PointF[T constraints.Float] struct {
	X, Y T
}

func BuildSVGPathDataF[T constraints.Float](points []PointF[T]) string {
	if len(points) < 3 {
		return ""
	}
	coords := make([]string, len(points))
	for i, p := range points {
		coords[i] = fmt.Sprintf("%g %g", p.X, p.Y)
	}
	return "M" + strings.Join(coords[:2], " L") + " " + strings.Join(coords[2:], " ") + " z"
}

func BuildSVGPathData(points []image.Point) string {

	if len(points) < 3 {
		return ""
	}

	coords := make([]string, len(points))
	for i := range points {
		coords[i] = fmt.Sprintf("%d %d", points[i].X, points[i].Y)
	}
	return "M" + strings.Join(coords[:2], " L") + " " + strings.Join(coords[2:], " ") + " z"

}
