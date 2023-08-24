package tag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (r Range) Len() int {
	return r.End - r.Begin
}

func (r Range) String() string {
	if r.Len() == 1 {
		return fmt.Sprintf("%d", r.Begin)
	}
	if r.End < 0 {
		return fmt.Sprintf("%d-", r.Begin)
	}
	if r.Begin <= 0 {
		return fmt.Sprintf("-%d", r.End)
	}

	return fmt.Sprintf("%d-%d", r.Begin, r.End)
}

type Range struct {
	Begin, End int
}

func (r *Range) UnmarshalFlag(s string) error {
	bounds := strings.Split(s, "-")
	if len(bounds) > 2 {
		return fmt.Errorf("Only supports ranges XX XX- -XX XX-YY, got '%s'", s)
	}
	if len(bounds) == 1 {
		idx, err := strconv.ParseInt(bounds[0], 10, 64)
		if err != nil {
			return fmt.Errorf("parsing range '%s': %w", s, err)
		}
		r.Begin = int(idx)
		r.End = r.Begin + 1
		return nil
	}

	if len(bounds[0]) == 0 {
		r.Begin = 0
	} else {
		idx, err := strconv.ParseInt(bounds[0], 10, 64)
		if err != nil {
			return fmt.Errorf("parsing range '%s': %w", s, err)
		}
		r.Begin = int(idx)
	}
	if len(bounds[1]) == 0 {
		r.End = -1
	} else {
		idx, err := strconv.ParseInt(bounds[1], 10, 64)
		if err != nil {
			return fmt.Errorf("parsing range '%s': %w", s, err)
		}
		r.End = int(idx)
	}
	return nil
}

func unmarshalRanges(s string) ([]Range, error) {
	ranges := strings.Split(s, ";")
	res := make([]Range, len(ranges))
	for i := range ranges {
		if err := res[i].UnmarshalFlag(ranges[i]); err != nil {
			return nil, fmt.Errorf("parsing ranges '%s': %w", s, err)
		}
	}
	return res, nil
}

type FamilyBlock struct {
	Family *Family
	SizeMM float64
	Ranges []Range
}

func (c *FamilyBlock) UnmarshalFlag(s string) error {
	wrapError := func(err error) error {
		return fmt.Errorf("invalid family specification '%s': %w", s, err)
	}

	args := strings.Split(s, ":")
	if len(args) > 3 || len(args) < 2 {
		return wrapError(errors.New("format should be <name>:<size>[:<range>]"))
	}

	var err error
	c.Family, err = GetFamily(args[0])
	if err != nil {
		return wrapError(err)
	}
	c.SizeMM, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return wrapError(err)
	}
	if len(args) == 2 {
		c.Ranges = []Range{{Begin: 0, End: len(c.Family.Codes)}}
		return nil
	}
	c.Ranges, err = unmarshalRanges(args[2])
	if err != nil {
		return wrapError(err)
	}
	return nil
}

func (c FamilyBlock) String() string {
	familyAndSize := fmt.Sprintf("%s:%g", c.Family.Name, c.SizeMM)
	if len(c.Ranges) == 1 && c.Ranges[0].Begin == 0 && c.Ranges[0].End == len(c.Family.Codes) {
		return familyAndSize
	}

	ranges := make([]string, len(c.Ranges))
	for i, r := range c.Ranges {
		ranges[i] = r.String()
	}

	return familyAndSize + ":" + strings.Join(ranges, ";")
}
