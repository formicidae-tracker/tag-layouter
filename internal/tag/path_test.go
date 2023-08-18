package tag

import (
	"fmt"
	"image"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PathSuite struct {
	suite.Suite
}

func TestPathSuite(t *testing.T) {
	suite.Run(t, new(PathSuite))
}

func renderImage(img *image.Gray) string {
	res := "\n┌" + strings.Repeat("─", img.Rect.Dx()) + "┐"
	for y := 0; y < img.Rect.Dy(); y++ {
		res += "\n│"
		for x := 0; x < img.Rect.Dx(); x++ {
			v := img.GrayAt(x, y).Y
			if v == 0x00 {
				res += "█"
			} else {
				res += " "
			}
		}
		res += "│"
	}
	return res + "\n└" + strings.Repeat("─", img.Rect.Dx()) + "┘\n"
}

func (s *PathSuite) requireGrayImage(img []uint8) *image.Gray {
	size := int(math.Sqrt(float64(len(img))))
	s.Require().Equal(size*size, len(img), "invalid test data: image must be square")
	return &image.Gray{
		Pix:    img,
		Stride: size,
		Rect:   image.Rect(0, 0, size, size),
	}
}

func (s *PathSuite) testPathBuilding(img []uint8, expected [][]image.Point) {
	gray := s.requireGrayImage(img)
	path := BuildPath(gray)
	s.Equal(expected, path)
}

func (s *PathSuite) TestSimplePolygonGeneration() {
	testdata := []struct {
		Name     string
		Image    []uint8
		Expected [][]image.Point
	}{
		{
			Name: "1x1 black",
			Image: []uint8{
				0x00,
			},
			Expected: [][]image.Point{
				{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			},
		},

		{
			Name: "3x3 black dot",
			Image: []uint8{
				0xff, 0xff, 0xff,
				0xff, 0x00, 0xff,
				0xff, 0xff, 0xff,
			},
			Expected: [][]image.Point{
				{{1, 1}, {2, 1}, {2, 2}, {1, 2}},
			},
		},
		{
			Name: "3x3 cross",
			Image: []uint8{
				0xff, 0x00, 0xff,
				0x00, 0x00, 0x00,
				0xff, 0x00, 0xff,
			},
			Expected: [][]image.Point{
				{
					{1, 0}, {2, 0}, {2, 1}, {3, 1}, {3, 2}, {2, 2},
					{2, 3}, {1, 3}, {1, 2}, {0, 2}, {0, 1}, {1, 1},
				},
			},
		},
		{
			Name: "3x3 diagonal",
			Image: []uint8{
				0x00, 0xff, 0xff,
				0xff, 0x00, 0xff,
				0xff, 0xff, 0x00,
			},
			Expected: [][]image.Point{
				{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
				{{1, 1}, {2, 1}, {2, 2}, {1, 2}},
				{{2, 2}, {3, 2}, {3, 3}, {2, 3}},
			},
		},
	}

	for _, d := range testdata {
		if s.Run(d.Name, func() { s.testPathBuilding(d.Image, d.Expected) }) == false {
			img := s.requireGrayImage(d.Image)
			fmt.Printf("failed image for '%s' is: %s", d.Name, renderImage(img))
		}
	}

}

func (s *PathSuite) TestComplexPolygonGeneration() {
	testdata := []struct {
		Name     string
		Image    []uint8
		Expected [][]image.Point
	}{
		{
			Name: "3x3 white dot",
			Image: []uint8{
				0x00, 0x00, 0x00,
				0x00, 0xff, 0x00,
				0x00, 0x00, 0x00,
			},
			Expected: [][]image.Point{
				{{0, 0}, {3, 0}, {3, 3}, {0, 3}},
				{{1, 1}, {2, 1}, {2, 2}, {1, 2}},
			},
		},
	}

	for _, d := range testdata {
		if s.Run(d.Name, func() { s.testPathBuilding(d.Image, d.Expected) }) == false {
			img := s.requireGrayImage(d.Image)
			fmt.Printf("failed image for '%s' is: %s", d.Name, renderImage(img))
		}
	}

}
