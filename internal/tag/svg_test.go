package tag

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SVGSuite struct {
	suite.Suite
}

func TestSVGSuite(t *testing.T) {
	suite.Run(t, new(SVGSuite))
}

func (s *SVGSuite) TestColorRendering() {
	testdata := []struct {
		Color    color.Color
		Expected string
	}{
		{color.Black, "#000000"},
		{color.White, "#ffffff"},
		{color.RGBA{127, 127, 127, 255}, "#7f7f7f"},
		{color.RGBA{0xff, 0x00, 0x00, 0xf0}, "#ff0000f0"},
	}

	for _, d := range testdata {
		s.Run(d.Expected, func() {
			s.Equal(toHex(d.Color), d.Expected)
		})
	}
}
