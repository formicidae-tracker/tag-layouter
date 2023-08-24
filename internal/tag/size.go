package tag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Size struct {
	Width, Height float64
}

func (s Size) String() string {
	return fmt.Sprintf("%gx%g", s.Width, s.Height)
}

func (s *Size) UnmarshalFlag(str string) error {
	wrapError := func(err error) error {
		return fmt.Errorf("invalid size specification '%s': %w", str, err)
	}

	dims := strings.Split(str, "x")

	if len(dims) != 2 {
		return wrapError(errors.New("format needs to be <width>x<height>"))
	}

	var err error
	s.Width, err = strconv.ParseFloat(dims[0], 64)
	if err != nil {
		return wrapError(err)
	}

	s.Height, err = strconv.ParseFloat(dims[1], 64)
	if err != nil {
		return wrapError(err)
	}

	return nil
}
