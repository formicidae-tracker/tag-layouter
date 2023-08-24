package tag

import (
	"fmt"
	"image/color"
	"log"

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
	log.Println(toHex(polygons[0].Color))
}
