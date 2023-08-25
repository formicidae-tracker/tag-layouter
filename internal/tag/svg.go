package tag

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	svg "github.com/ajstarks/svgo"
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
		paths[i] = buildSVGPath(polygons[i].Vertices)
	}

	SVG.Path(strings.Join(paths, " "),
		fmt.Sprintf("style=\"fill:%s;\"", toHex(polygons[0].Color)))
}

type PointF[T float32 | float64] struct {
	X, Y T
}

func BuildSVGD[T float32 | float64](points []PointF[T]) string {
	if len(points) < 2 {
		return ""
	}
	coords := make([]string, len(points))
	for i, p := range points {
		coords[i] = fmt.Sprintf("%g,%g", p.X, p.Y)
	}
	return "M " + strings.Join(coords, " L ") + " z"
}

func buildSVGPath(points []image.Point) string {

	if len(points) < 2 {
		return ""
	}

	coords := make([]string, len(points))
	for i := range points {
		coords[i] = fmt.Sprintf("%d,%d", points[i].X, points[i].Y)
	}
	return "M " + strings.Join(coords, " L ") + " z"

}
