package tag

import "image"

type Polygon []image.Point

type Tag struct {
	Image   image.Gray
	Polygon []Polygon
}
