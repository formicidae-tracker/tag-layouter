package main

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
	"strings"
	"unsafe"
)

type TagFamily struct {
	Codes          []uint64
	NBits          int
	LocationX      []int
	LocationY      []int
	Name           string
	Hamming        int
	TotalWidth     int
	WidthAtBorder  int
	ReversedBorder bool
}

func newTagFamily(tf *C.apriltag_family_t, name string) *TagFamily {
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

	res := &TagFamily{
		Codes:          append([]uint64{}, codes...),
		Name:           name,
		NBits:          nbits,
		TotalWidth:     int(tf.total_width),
		WidthAtBorder:  int(tf.width_at_border),
		ReversedBorder: bool(tf.reversed_border),
		Hamming:        int(tf.h),
		LocationX:      nil,
		LocationY:      nil,
	}
	for i := 0; i < res.NBits; i++ {
		res.LocationX = append(res.LocationX, int(bitX[i]))
		res.LocationY = append(res.LocationY, int(bitY[i]))
	}
	return res
}

func (tf *TagFamily) CodeSize() int {
	return tf.NBits
}

func (tf *TagFamily) SizeInPX() int {
	return tf.TotalWidth
}

type cAprilTagFamily struct {
	Constructor func() *C.apriltag_family_t
	Destructor  func(*C.apriltag_family_t)
}

var familyFactory map[string]cAprilTagFamily

func init() {
	familyFactory = map[string]cAprilTagFamily{
		"36h10": cAprilTagFamily{
			Constructor: func() *C.apriltag_family_t { return C.tag36h10_create() },
			Destructor:  func(f *C.apriltag_family_t) { C.tag36h10_destroy(f) },
		},
		"36h11": cAprilTagFamily{
			Constructor: func() *C.apriltag_family_t { return C.tag36h11_create() },
			Destructor:  func(f *C.apriltag_family_t) { C.tag36h11_destroy(f) },
		},
		"16h5": cAprilTagFamily{
			Constructor: func() *C.apriltag_family_t { return C.tag16h5_create() },
			Destructor:  func(f *C.apriltag_family_t) { C.tag16h5_destroy(f) },
		},
		"25h9": cAprilTagFamily{
			Constructor: func() *C.apriltag_family_t { return C.tag25h9_create() },
			Destructor:  func(f *C.apriltag_family_t) { C.tag25h9_destroy(f) },
		},
		"Circle21h7": cAprilTagFamily{
			Constructor: func() *C.apriltag_family_t { return C.tagCircle21h7_create() },
			Destructor:  func(f *C.apriltag_family_t) { C.tagCircle21h7_destroy(f) },
		},
		"Circle49h12": cAprilTagFamily{
			Constructor: func() *C.apriltag_family_t { return C.tagCircle49h12_create() },
			Destructor:  func(f *C.apriltag_family_t) { C.tagCircle49h12_destroy(f) },
		},
		"Custom48h12": cAprilTagFamily{
			Constructor: func() *C.apriltag_family_t { return C.tagCustom48h12_create() },
			Destructor:  func(f *C.apriltag_family_t) { C.tagCustom48h12_destroy(f) },
		},
		"Standard41h12": cAprilTagFamily{
			Constructor: func() *C.apriltag_family_t { return C.tagStandard41h12_create() },
			Destructor:  func(f *C.apriltag_family_t) { C.tagStandard41h12_destroy(f) },
		},
		"Standard52h13": cAprilTagFamily{
			Constructor: func() *C.apriltag_family_t { return C.tagStandard52h13_create() },
			Destructor:  func(f *C.apriltag_family_t) { C.tagStandard52h13_destroy(f) },
		},
	}
}

func GetFamily(name string) (*TagFamily, error) {
	ff, ok := familyFactory[name]
	if ok == false {
		return nil, fmt.Errorf("Unknown famnily '%s'", name)
	}
	name = strings.ToUpper(name)
	tf := ff.Constructor()
	defer ff.Destructor(tf)
	return newTagFamily(tf, name), nil
}
