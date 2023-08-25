package tag

import (
	"math"

	"golang.org/x/exp/constraints"
)

const anInchInMM float64 = 25.4

func PixelToMM[T constraints.Float | constraints.Integer](DPI int, v T) float64 {
	return float64(v) * anInchInMM / float64(DPI)
}

func MMToPixel[T constraints.Float | constraints.Integer](DPI int, v T) int {
	return int(math.Round(float64(v) * float64(DPI) / anInchInMM))
}
