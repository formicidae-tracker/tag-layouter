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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"unsafe"

	"github.com/formicidae-tracker/tag-layouter/internal/tag"
	"github.com/jessevdk/go-flags"
	"golang.org/x/exp/slog"
)

func newTagFamily(tf *C.apriltag_family_t, name string) *tag.Family {
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

	res := &tag.Family{
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

type Options struct {
	Args struct {
		OutputDir flags.Filename
	} `positional-args:"yes" required:"yes"`
}

func main() {
	if err := execute(); err != nil {
		log.Fatalf("unhandled error: %s", err)
	}
}

func execute() error {
	var opts Options
	if _, err := flags.Parse(&opts); err != nil {
		if flags.WroteHelp(err) == true {
			os.Exit(0)
		}
		os.Exit(1)
	}

	os.MkdirAll(string(opts.Args.OutputDir), 0755)

	for name, factory := range familyFactory {
		if err := opts.generateJSON(name, factory); err != nil {
			return fmt.Errorf("generating %s: %w", name, err)
		}
	}

	return nil
}

func (o Options) generateJSON(name string, factory cAprilTagFamily) error {
	filename := filepath.Join(string(o.Args.OutputDir), name+".json")
	slog.Info("generating", "family", name, "filepath", filename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	family := newTagFamily(factory(), name)
	enc := json.NewEncoder(f)
	return enc.Encode(family)
}
