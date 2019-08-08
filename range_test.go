package main

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type RangeSuite struct{}

var _ = Suite(&RangeSuite{})

func (s *RangeSuite) TestExtractRange(c *C) {
	testdata := []struct {
		Input    string
		Expected Range
	}{
		{
			"42",
			Range{42, 43},
		},
		{
			"-42",
			Range{0, 42},
		},
		{
			"42-",
			Range{42, -1},
		},
	}

	for _, d := range testdata {
		r, err := ExtractRange(d.Input)
		c.Check(err, IsNil, Commentf("Parsing '%s'", d.Input))
		c.Check(r, Equals, d.Expected)
	}

}

func (s *RangeSuite) TestExtractRanges(c *C) {
	testdata := []struct {
		Input    string
		Expected []Range
	}{
		{
			"",
			nil,
		},
		{
			"0;1;2;3",
			[]Range{Range{0, 1}, Range{1, 2}, Range{2, 3}, Range{3, 4}},
		},
		{
			"-41;42-",
			[]Range{Range{0, 41}, Range{42, -1}},
		},
	}

	for _, d := range testdata {
		r, err := ExtractRanges(d.Input)
		c.Check(err, IsNil, Commentf("Parsing '%s'", d.Input))
		c.Check(r, DeepEquals, d.Expected)
	}

}
