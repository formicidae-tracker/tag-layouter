package tag

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FamilyBlockSuite struct {
	suite.Suite
}

func TestFamilyBlockSuite(t *testing.T) {
	suite.Run(t, new(FamilyBlockSuite))
}

func (s *FamilyBlockSuite) TestRangeFlagIO() {
	testdata := []struct {
		R        Range
		Expected string
	}{
		{Range{0, -1}, "0-"},
		{Range{3, -1}, "3-"},
		{Range{3, 4}, "3"},
		{Range{0, 4}, "-4"},
		{Range{1, 4}, "1-4"},
	}

	for _, d := range testdata {
		s.Run(d.Expected, func() {
			s.Equal(d.Expected, d.R.String(), "formatting")
			var res Range
			s.Require().NoError(res.UnmarshalFlag(d.Expected))
			s.Equal(d.R, res, "parsing")
		})
	}
}

func (s *FamilyBlockSuite) TestFamilyConfigIO() {
	mustFamily := func(name string) *Family {
		f, err := GetFamily(name)
		if err != nil {
			panic(err.Error())
		}
		return f
	}
	testdata := []struct {
		C        FamilyBlock
		Expected string
	}{
		{FamilyBlock{mustFamily("36h11"), 1.6, []Range{{0, 587}}}, "36h11:1.6"},
		{FamilyBlock{mustFamily("36h11"), 1.6, []Range{{1, 5}, {8, -1}}}, "36h11:1.6:1-5;8-"},
	}

	for _, d := range testdata {
		s.Run(d.Expected, func() {
			s.Equal(d.Expected, d.C.String(), "formatting")
			var res FamilyBlock
			s.Require().NoError(res.UnmarshalFlag(d.Expected))
			s.Equal(d.C, res, "parsing")
		})
	}
}
