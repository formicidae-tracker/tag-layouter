package main

type Dotter struct {
	dpi float64
}

const anInch = 25.4

func (d Dotter) ToDot(v float64) int {
	return int(v * d.dpi / anInch)
}

func (d Dotter) ToMM(v int) float64 {
	return float64(v) * anInch / float64(d.dpi)
}
