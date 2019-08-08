package main

import "fmt"

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

func ExtractRanges(s string) ([]Range, error) {
	// 	ranges := strings.Split(fargs[2], "-")
	// 	begin := -1
	// 	end := -1
	// 	if len(ranges) > 2 {
	// 		return res, fmt.Errorf("Only supports ranges XX XX- -XX XX-YY, got '%s'", fargs[2])
	// 	}
	// 	if len(ranges) == 1 {
	// 		idx, err := strconv.ParseInt(ranges[0], 10, 64)
	// 		if err != nil {
	// 			return res, err
	// 		}
	// 		begin = int(idx)
	// 		end = int(idx) + 1

	// 	} else {
	// 		if len(ranges[0]) == 0 {
	// 			begin = 0
	// 		} else {
	// 			idx, err := strconv.ParseInt(ranges[0], 10, 64)
	// 			if err != nil {
	// 				return res, err
	// 			}
	// 			if int(idx) >= len(tf.Codes) {
	// 				return res, fmt.Errorf("%d is out-of-range in %s (size:%d)'", idx, fargs[0], len(tf.Codes))
	// 			}
	// 			begin = int(idx)
	// 		}
	// 		if len(ranges[1]) == 0 {
	// 			end = len(tf.Codes)
	// 		} else {
	// 			idx, err := strconv.ParseInt(ranges[0], 10, 64)
	// 			if err != nil {
	// 				return res, err
	// 			}
	// 			if int(idx) >= len(tf.Codes) {
	// 				return res, fmt.Errorf("%d is out-of-range in %s (size:%d)'", idx, fargs[0], len(tf.Codes))
	// 			}
	// 			end = int(idx)
	// 		}
	// 	}
	// 	res = append(res, FamilyAndSize{
	// 		Family: tf,
	// 		Size:   s,
	// 		Begin:  begin,
	// 		End:    end,
	// 	})
	// }
	return nil, fmt.Errorf("not yet implemented")
}

type FamilyBlock struct {
	Family *TagFamily
	Size   float64
	Ranges []Range
}

func (f *FamilyBlock) FamilyLabel() string {
	return fmt.Sprintf("%s %.2fMM", f.Family.Name, f.Size)
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
