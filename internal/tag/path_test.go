package tag

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/image/vector"
)

type PathSuite struct {
	suite.Suite
}

func TestPathSuite(t *testing.T) {
	suite.Run(t, new(PathSuite))
}

func renderImage(img *image.Gray) string {
	res := "\n┌" + strings.Repeat("─", 2*img.Rect.Dx()) + "┐"
	for y := 0; y < img.Rect.Dy(); y++ {
		res += "\n│"
		for x := 0; x < img.Rect.Dx(); x++ {
			v := img.GrayAt(x, y).Y
			if v == 0xff {
				res += "██"
			} else if v != 0x00 {
				res += "xx"
			} else {
				res += "  "
			}
		}
		res += "│"
	}
	return res + "\n└" + strings.Repeat("─", 2*img.Rect.Dx()) + "┘\n"
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
		{
			Name: "36h11-00 1st blob",
			Image: []uint8{
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0x00, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff,
				0xff, 0xff, 0x00, 0x00, 0x00, 0xff, 0x00, 0xff,
				0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
				0xff, 0x00, 0xff, 0x00, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0x00, 0xff, 0x00, 0x00, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0x00, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
			Expected: [][]image.Point{
				{{1, 1}, {3, 1}, {3, 2}, {4, 2}, {4, 1}, {5, 1}, {5, 3}, {4, 3},
					{4, 5}, {3, 5}, {3, 4}, {2, 4}, {2, 2}, {1, 2}},
				{{6, 1}, {7, 1}, {7, 3}, {6, 3}},
				{{1, 4}, {2, 4}, {2, 5}, {1, 5}},
				{{2, 5}, {3, 5}, {3, 6}, {2, 6}},
				{{4, 5}, {6, 5}, {6, 6}, {5, 6}, {5, 7}, {4, 7}},
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

func (s *PathSuite) TestTagRendering() {
	families := []string{"36h11", "Standard41h12"}

	for _, f := range families {
		s.Run(f, func() { s.testTagFamily(f) })
	}
}

func (s *PathSuite) testTagFamily(name string) {
	family, err := GetFamily(name)
	s.Require().NoError(err)
	for i := range family.Codes {
		if s.Run(fmt.Sprintf("Tag %08d", i), func() {
			s.testTagRendering(family, i)
		}) == false {
			return
		}
	}
}

func (s *PathSuite) testTagRendering(family *Family, i int) {
	img := family.BuildTag(i)
	paths := BuildPath(img)
	colorInversions := 1
	if family.ReversedBorder {
		colorInversions = 2
	}

	rendered := s.renderTagFromPath(paths, img.Rect, colorInversions)
	if s.Equal(img, rendered) == false {
		renderImageDiff(os.Stdout, i, img, rendered)
	}
}

func (s *PathSuite) renderTagFromPath(paths [][]image.Point, size image.Rectangle, colorInversions int) *image.Gray {
	s.Require().GreaterOrEqual(len(paths), colorInversions+1)

	//start with a white image
	dst := image.NewGray(size)
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	for i := 0; i < colorInversions; i++ {
		penColor := color.Black
		if i%2 == 1 {
			penColor = color.White
		}
		s.drawPathOn(paths[i], dst, penColor)
	}

	penColor := color.Black
	if colorInversions%2 == 1 {
		penColor = color.White
	}
	for _, path := range paths[colorInversions:] {
		s.drawPathOn(path, dst, penColor)
	}

	return dst
}

func (s *PathSuite) drawPathOn(path []image.Point, img *image.Gray, c color.Color) {
	s.Require().Greater(len(path), 0)
	r := vector.NewRasterizer(img.Bounds().Dx(), img.Bounds().Dy())
	r.MoveTo(float32(path[0].X), float32(path[0].Y))
	for _, v := range path[1:] {
		r.LineTo(float32(v.X), float32(v.Y))
	}
	r.ClosePath()

	r.Draw(img, img.Bounds(), image.NewUniform(c), image.Point{})
}

func renderImageDiff(w io.Writer, i int, expected *image.Gray, actual *image.Gray) {
	if expected.Rect != actual.Rect {
		return
	}

	topBorder := strings.Repeat("─", 2*expected.Rect.Dx())

	titleFormat := fmt.Sprintf("%%-%ds %%04d", 2*expected.Rect.Dx()-5)
	fmt.Fprintf(w, " "+titleFormat+"   "+titleFormat+"   "+titleFormat+"\n", "Actual", i, "Expected", i, "Diff", i)
	fmt.Fprintf(w, "┌%s┐ ┌%s┐ ┌%s┐\n", topBorder, topBorder, topBorder)
	for y := 0; y < expected.Rect.Dy(); y++ {
		aLine := ""
		eLine := ""
		dLine := ""
		for x := 0; x < expected.Rect.Dx(); x++ {
			aV := actual.GrayAt(x, y).Y
			eV := expected.GrayAt(x, y).Y
			if aV == 0x00 {
				aLine += "██"
			} else {
				aLine += "  "
			}
			if eV == 0x00 {
				eLine += "██"
				if aV == eV {
					dLine += "██"
				} else {
					dLine += "XX"
				}
			} else {
				eLine += "  "
				if aV == eV {
					dLine += "  "
				} else {
					dLine += "··"
				}
			}
		}
		fmt.Fprintf(w, "│%s│ │%s│ │%s│\n", aLine, eLine, dLine)

	}

	fmt.Fprintf(w, "└%s┘ └%s┘ └%s┘\n", topBorder, topBorder, topBorder)

}
