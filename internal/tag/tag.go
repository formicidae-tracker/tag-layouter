package tag

type Image struct {
	Height int
	Width  int
	Pixels []uint8
}

func newImage(w, h int) Image {
	return Image{
		Height: h,
		Width:  w,
		Pixels: make([]uint8, w*h),
	}
}

func (i *Image) Set(x, y int, value uint8) {
	i.Pixels[y*i.Width+x] = value
}

func (i *Image) Get(x, y int) uint8 {
	return i.Pixels[y*i.Width+x]
}

func (i Image) Map(fn func(x, y int, value uint8)) {
	for k, v := range i.Pixels {
		y := k / i.Width
		x := k - y*i.Width
		fn(x, y, v)
	}
}

type Vertex struct {
	X, Y int
}

type Polygon []Vertex

type Tag struct {
	Image   Image
	Polygon []Polygon
}
