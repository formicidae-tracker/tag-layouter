package tag

/*
#include "apriltag/apriltag.h"
#include "apriltag/tag16h5.h"
#include "apriltag/tag25h9.h"
#include "apriltag/tag36h11.h"
#include "apriltag/tagCircle21h7.h"
#include "apriltag/tagCircle49h12.h"
#include "apriltag/tagCustom48h12.h"
#include "apriltag/tagStandard41h12.h"
#include "apriltag/tagStandard52h13.h"
#include "oldtags/tag36h10.h"
#cgo LDFLAGS:  apriltag/libapriltag.a oldtags/liboldtags.a -lm
*/
import "C"

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"unsafe"
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

func newTagFamily(tf *C.apriltag_family_t, name string) *Family {
	ncodes := int(tf.ncodes)
	nbits := int(tf.nbits)
	var codes []uint64
	var bitX, bitY []int32

	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&codes)))
	sliceHeader.Cap = ncodes
	sliceHeader.Len = ncodes
	sliceHeader.Data = uintptr(unsafe.Pointer(tf.codes))
	sliceHeaderX := (*reflect.SliceHeader)((unsafe.Pointer(&bitX)))
	sliceHeaderX.Cap = nbits
	sliceHeaderX.Len = nbits
	sliceHeaderX.Data = uintptr(unsafe.Pointer(tf.bit_x))
	sliceHeaderY := (*reflect.SliceHeader)((unsafe.Pointer(&bitY)))
	sliceHeaderY.Cap = nbits
	sliceHeaderY.Len = nbits
	sliceHeaderY.Data = uintptr(unsafe.Pointer(tf.bit_y))

	res := &Family{
		Codes:          append([]uint64{}, codes...),
		Name:           name,
		NBits:          nbits,
		TotalWidth:     int(tf.total_width),
		WidthAtBorder:  int(tf.width_at_border),
		ReversedBorder: bool(tf.reversed_border),
		Hamming:        int(tf.h),
		LocationX:      nil,
		LocationY:      nil,
		Inside:         nil,
	}
	for i := 0; i < res.NBits; i++ {
		x := int(bitX[i])
		y := int(bitY[i])
		res.LocationX = append(res.LocationX, x)
		res.LocationY = append(res.LocationY, y)
		inside := true
		if x < 0 || x >= res.WidthAtBorder {
			inside = false
		}
		if y < 0 || y >= res.WidthAtBorder {
			inside = false
		}
		res.Inside = append(res.Inside, inside)
	}

	return res
}

func (tf *Family) CodeSize() int {
	return tf.NBits
}

func (tf *Family) SizeInPX() int {
	return tf.TotalWidth
}

type cAprilTagFamily func() *C.apriltag_family_t

var familyFactory map[string]cAprilTagFamily

func init() {
	familyFactory = map[string]cAprilTagFamily{
		"36h10":         func() *C.apriltag_family_t { return C.tag36h10_create() },
		"36h11":         func() *C.apriltag_family_t { return C.tag36h11_create() },
		"16h5":          func() *C.apriltag_family_t { return C.tag16h5_create() },
		"25h9":          func() *C.apriltag_family_t { return C.tag25h9_create() },
		"Circle21h7":    func() *C.apriltag_family_t { return C.tagCircle21h7_create() },
		"Circle49h12":   func() *C.apriltag_family_t { return C.tagCircle49h12_create() },
		"Custom48h12":   func() *C.apriltag_family_t { return C.tagCustom48h12_create() },
		"Standard41h12": func() *C.apriltag_family_t { return C.tagStandard41h12_create() },
		"Standard52h13": func() *C.apriltag_family_t { return C.tagStandard52h13_create() },
	}
}

var allocated = make(map[string]*Family)

// GetFamily, returns a Family given its canonical name. It may return
// an error if the canonical name is unknown.
func GetFamily(name string) (*Family, error) {
	if alreadyAllocated, ok := allocated[name]; ok == true {
		return alreadyAllocated, nil
	}

	ff, ok := familyFactory[name]
	if ok == false {
		return nil, fmt.Errorf("Unknown famnily '%s'", name)
	}
	res := newTagFamily(ff(), name)
	allocated[name] = res
	return res, nil
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
