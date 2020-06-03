package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Range struct {
	Begin, End int
}

func (r Range) Len() int {
	return r.End - r.Begin
}

func (r Range) String() string {
	if r.Len() == 1 {
		return fmt.Sprintf("%d", r.Begin)
	}
	return fmt.Sprintf("%d-%d", r.Begin, r.End)
}

func ExtractRange(s string) (Range, error) {
	ranges := strings.Split(s, "-")
	if len(ranges) > 2 {
		return Range{}, fmt.Errorf("Only supports ranges XX XX- -XX XX-YY, got '%s'", s)
	}
	if len(ranges) == 1 {
		idx, err := strconv.ParseInt(ranges[0], 10, 64)
		if err != nil {
			return Range{}, err
		}
		return Range{Begin: int(idx), End: int(idx + 1)}, nil

	}
	begin := -1
	end := -1
	if len(ranges[0]) == 0 {
		begin = 0
	} else {
		idx, err := strconv.ParseInt(ranges[0], 10, 64)
		if err != nil {
			return Range{}, err
		}
		begin = int(idx)
	}
	if len(ranges[1]) == 0 {
		end = -1
	} else {
		idx, err := strconv.ParseInt(ranges[1], 10, 64)
		if err != nil {
			return Range{}, err
		}
		end = int(idx)
	}
	return Range{Begin: begin, End: end}, nil
}

func ExtractRanges(s string) ([]Range, error) {
	rangesStr := strings.Split(s, ";")
	if len(rangesStr) == 1 && len(rangesStr[0]) == 0 {
		return nil, nil
	}
	var res []Range
	for _, rStr := range rangesStr {
		if len(rStr) == 0 {
			continue
		}
		r, err := ExtractRange(rStr)
		if err != nil {
			return res, err
		}
		res = append(res, r)
	}

	return res, nil
}

type FamilyBlock struct {
	Family *TagFamily
	Size   float64
	Ranges []Range
}

func (f *FamilyBlock) FamilyLabelActualSize(size float64) string {
	return fmt.Sprintf("%s %.2fMM", f.Family.Name, size)
}

func (f *FamilyBlock) FamilyLabel() string {
	return f.FamilyLabelActualSize(f.Size)
}

func (f *FamilyBlock) NumberOfTags() int {
	n := 0
	for _, r := range f.Ranges {
		n += r.Len()
	}
	return n
}

func (f *FamilyBlock) RangeString() string {
	res := ""
	sep := ""
	for _, r := range f.Ranges {
		res = fmt.Sprintf("%s%s%s", res, sep, r)
		sep = ";"
	}
	return res
}
