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
	"reflect"
	"unsafe"
)

type Family struct {
	Codes          []uint64
	NBits          int
	LocationX      []int
	LocationY      []int
	Inside         []bool
	Name           string
	Hamming        int
	TotalWidth     int
	WidthAtBorder  int
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
