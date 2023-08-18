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

func renderImage(img []uint8) string {
	size := int(math.Sqrt(float64(len(img))))
	res := "\n┌" + strings.Repeat("─", size) + "┐"
	for y := 0; y < size; y++ {
		res += "\n│"
		for x := 0; x < size; x++ {
			v := img[y*size+x]
			if v == 0x00 {
				res += "█"
			} else {
				res += " "
			}
		}
		res += "│"
	}
	return res + "\n└" + strings.Repeat("─", size) + "┘\n"
}

func (s *PathSuite) testPathBuilding(img []uint8, expected [][]image.Point) {
	size := int(math.Sqrt(float64(len(img))))
	s.Require().Equal(size*size, len(img), "invalid test data: image must be square")
	gray := &image.Gray{
		Pix:    img,
		Stride: size,
		Rect:   image.Rect(0, 0, size, size),
	}
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
				{{1, 1}, {2, 1}, {2, 2}, {1, 2}},
			},
		},
	}

	for _, d := range testdata {
		if s.Run(d.Name, func() { s.testPathBuilding(d.Image, d.Expected) }) == false {
			fmt.Printf("failed image for '%s' is: %s", d.Name, renderImage(d.Image))
		}
	}

}
