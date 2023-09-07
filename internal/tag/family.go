package tag

import (
	"embed"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"path/filepath"
	"strings"
)

// A Family defines how tags are defined: their size, the nuuber of
// available codes, and how it should be
// rendered. <Family.RenderTag>() can be used to produce monochromatic
// image.GrayImage that represent the tag.
type Family struct {
	// The list of codes available in the family.
	Codes []uint64
	// Number of coding bits in the family
	NBits int
	// X location of coding bit number i
	LocationX []int
	// Y location of coding bit number i
	LocationY []int
	// True if the bit is inside the border or outside the border
	Inside []bool
	// Canonical name of the family
	Name string
	//  Minimal Hamming distance between any two code in the family:
	//  The higher, the better.
	Hamming int
	// Length size of the total tag in pixel
	TotalWidth int
	// Length size at the white/black tag border.
	WidthAtBorder int
	// If true, the outside border is black and the inside border is
	// white.
	ReversedBorder bool
}

func (tf *Family) CodeSize() int {
	return tf.NBits
}

func (tf *Family) SizeInPX() int {
	return tf.TotalWidth
}

//go:embed data/*.json
var familyData embed.FS

var familyNames []string

// GetFamily, returns a Family given its canonical name. It may return
// an error if the canonical name is unknown.
func GetFamily(name string) (*Family, error) {
	f, err := familyData.Open(filepath.Join("data", name+".json"))
	if err != nil {
		return nil, fmt.Errorf("could not get family '%s': %w", name, err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	res := &Family{}

	return res, dec.Decode(res)
}

func init() {
	entries, err := familyData.ReadDir("data")
	if err != nil {
		panic(err.Error())
	}
	for _, e := range entries {
		familyNames = append(familyNames, strings.TrimSuffix(e.Name(), ".json"))
	}
}

// RenderTag renders in a monochromatic image.GrayImage a tag given its code index.
//
// The rendered image has a size of (.TotalWidth,.TotalWidth) in
// pixel, and pixel value are either 0xff (background value) or 0x00
// (foreground value).
func (f *Family) RenderTag(n int) *image.Gray {
	res := image.NewGray(image.Rect(0, 0, f.TotalWidth, f.TotalWidth))
	offset := (f.TotalWidth - f.WidthAtBorder) / 2
	inside := image.Rect(offset, offset, offset+f.WidthAtBorder, offset+f.WidthAtBorder)
	if f.ReversedBorder == false {
		draw.Draw(res, res.Bounds(), image.NewUniform(color.White), image.Pt(0, 0), draw.Src)
		draw.Draw(res, inside, image.NewUniform(color.Black), image.Pt(0, 0), draw.Src)
	} else {
		draw.Draw(res, inside, image.NewUniform(color.White), image.Pt(0, 0), draw.Src)
	}
	for i := 0; i < f.NBits; i++ {
		var bit uint64 = (uint64(1) << (uint64(f.NBits) - 1 - uint64(i)))
		isSet := bit&f.Codes[n] != 0
		backgroundIsBlack := f.Inside[i] != f.ReversedBorder

		if isSet != backgroundIsBlack {
			continue
		}
		value := color.White
		if backgroundIsBlack == false {
			value = color.Black
		}
		res.Set(offset+f.LocationX[i], offset+f.LocationY[i], value)
	}

	return res
}
