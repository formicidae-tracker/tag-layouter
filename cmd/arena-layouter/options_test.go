package main

import (
	"testing"

	"gihtub.com/formicidae-tracker/tag-layouter/internal/tag"
	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/suite"
)

type OptionsSuite struct {
	suite.Suite
}

func TestOptionsSuite(t *testing.T) {
	suite.Run(t, new(OptionsSuite))
}

func (s *OptionsSuite) TestValidation() {
	testdata := []struct {
		O        Options
		Expected string
	}{
		{Options{Args: struct{ File flags.Filename }{File: "foo.png"}}, "invalid filepath 'foo.png': only SVG are supported, filepath must end with '.svg'"},

		{
			Options{
				Args:      struct{ File flags.Filename }{File: "good.svg"},
				ArenaSize: tag.Size{20, 10},
				PaperSize: tag.Size{10, 10},
			},
			"incompatible paper size (10x10) and arena size (20x10): the arena must fit on the paper",
		},
		{
			Options{
				Args:      struct{ File flags.Filename }{File: "good.svg"},
				ArenaSize: tag.Size{10, 20},
				PaperSize: tag.Size{10, 10},
			},
			"incompatible paper size (10x10) and arena size (10x20): the arena must fit on the paper",
		},
		{
			Options{
				Args:      struct{ File flags.Filename }{File: "good.svg"},
				ArenaSize: tag.Size{10, 10},
				PaperSize: tag.Size{20.3, 20},
			},
			"invalid paper size '20.3x20': sub-millimeter paper size are not supported",
		},
		{
			Options{
				Args:      struct{ File flags.Filename }{File: "good.svg"},
				ArenaSize: tag.Size{10, 10},
				PaperSize: tag.Size{20, 20.2},
			},
			"invalid paper size '20x20.2': sub-millimeter paper size are not supported",
		},
	}

	for _, d := range testdata {
		s.Run(d.Expected, func() {
			s.ErrorContains(d.O.Validate(), d.Expected)
		})
	}
}
