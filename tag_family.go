package main

/*
#include "apriltag-2016-12-01/apriltag.h"
#include "apriltag-2016-12-01/tag16h5.h"
#include "apriltag-2016-12-01/tag25h7.h"
#include "apriltag-2016-12-01/tag25h9.h"
#include "apriltag-2016-12-01/tag36h10.h"
#include "apriltag-2016-12-01/tag36h11.h"
#include "apriltag-2016-12-01/tag36artoolkit.h"
#cgo LDFLAGS:  apriltag-2016-12-01/libapriltag.a -lm
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

type TagFamily struct {
	Codes       []int64
	Name        string
	Size        int
	Hamming     int
	BlackBorder int
	WhiteBorder int
}

func newTagFamily(tf *C.apriltag_family_t, name string) *TagFamily {
	ncodes := int(tf.ncodes)
	var codes []int64

	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&codes)))
	sliceHeader.Cap = ncodes
	sliceHeader.Len = ncodes
	sliceHeader.Data = uintptr(unsafe.Pointer(tf.codes))

	return &TagFamily{
		Codes:       append([]int64{}, codes...),
		Name:        name,
		Size:        int(tf.d),
		WhiteBorder: 1,
		BlackBorder: int(tf.black_border),
		Hamming:     int(tf.d),
	}
}

func (tf *TagFamily) CodeSize() int {
	return tf.Size * tf.Size
}

func (tf *TagFamily) SizeInPX() int {
	return tf.Size + 2*(tf.WhiteBorder+tf.BlackBorder)
}

func GetFamily(name string) (*TagFamily, error) {
	if name == "36h11" {
		tf := C.tag36h11_create()
		defer C.tag36h11_destroy(tf)
		return newTagFamily(tf, "36H11"), nil
	} else if name == "36h10" {
		tf := C.tag36h10_create()
		defer C.tag36h10_destroy(tf)
		return newTagFamily(tf, "36H10"), nil
	} else if name == "25h9" {
		tf := C.tag25h9_create()
		defer C.tag25h9_destroy(tf)
		return newTagFamily(tf, "25H9"), nil
	} else if name == "25h7" {
		tf := C.tag25h7_create()
		defer C.tag25h7_destroy(tf)
		return newTagFamily(tf, "25H7"), nil
	} else if name == "16h5" {
		tf := C.tag16h5_create()
		defer C.tag16h5_destroy(tf)
		return newTagFamily(tf, "16H5"), nil
	} else if name == "36artoolkit" {
		tf := C.tag36artoolkit_create()
		defer C.tag36artoolkit_destroy(tf)
		return newTagFamily(tf, "ARTOOLKIT"), nil
	} else {
		return nil, fmt.Errorf("Unknown family %s", name)
	}
}
