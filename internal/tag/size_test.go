package tag

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SizeSuite struct {
	suite.Suite
}

func TestSizeSuite(t *testing.T) {
	suite.Run(t, new(SizeSuite))
}

func (s *SizeSuite) TestFormatIO() {
	testdata := []struct {
		Size     Size
		Expected string
	}{
		{Size{210, 297}, "210x297"},
		{Size{10.5, 12.8}, "10.5x12.8"},
	}

	for _, d := range testdata {
		s.Run(d.Expected, func() {
			s.Equal(d.Expected, d.Size.String())
			var res Size
			s.Require().NoError(res.UnmarshalFlag(d.Expected))
			s.Equal(d.Size, res)
		})
	}
}
