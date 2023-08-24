package tag

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FamilySuite struct {
	suite.Suite
}

func (s *FamilySuite) TestTagFamily() {
	testData := map[string]struct {
		NCodes int
		Family Family
	}{

		"16h5": {
			NCodes: 30,
			Family: Family{
				Codes:          []uint64{0x27c8},
				NBits:          16,
				Name:           "16h5",
				Hamming:        5,
				TotalWidth:     8,
				WidthAtBorder:  6,
				ReversedBorder: false,
			},
		},

		"25h9": {
			NCodes: 35,
			Family: Family{
				Codes:          []uint64{0x156f1f4},
				NBits:          25,
				Name:           "25h9",
				Hamming:        9,
				TotalWidth:     9,
				WidthAtBorder:  7,
				ReversedBorder: false,
			},
		},

		"36h10": {
			NCodes: 2320,
			Family: Family{
				Codes:          []uint64{0x1a42f9469},
				NBits:          36,
				Name:           "36h10",
				Hamming:        10,
				TotalWidth:     10,
				WidthAtBorder:  8,
				ReversedBorder: false,
			},
		},
		"36h11": {
			NCodes: 587,
			Family: Family{
				Codes:          []uint64{0xd7e00984b},
				NBits:          36,
				Name:           "36h11",
				Hamming:        11,
				TotalWidth:     10,
				WidthAtBorder:  8,
				ReversedBorder: false,
			},
		},

		"Circle21h7": {
			NCodes: 38,
			Family: Family{
				Codes:          []uint64{0x157863},
				NBits:          21,
				Name:           "Circle21h7",
				Hamming:        7,
				TotalWidth:     9,
				WidthAtBorder:  5,
				ReversedBorder: true,
			},
		},

		"Circle49h12": {
			NCodes: 65535,
			Family: Family{
				Codes:          []uint64{0xc6c921d8614a},
				NBits:          49,
				Name:           "Circle49h12",
				Hamming:        12,
				TotalWidth:     11,
				WidthAtBorder:  5,
				ReversedBorder: true,
			},
		},

		"Custom48h12": {
			NCodes: 42211,
			Family: Family{
				Codes:          []uint64{0xd6c8ae76dff0},
				NBits:          48,
				Name:           "Custom48h12",
				Hamming:        12,
				TotalWidth:     10,
				WidthAtBorder:  6,
				ReversedBorder: true,
			},
		},

		"Standard41h12": {
			NCodes: 2115,
			Family: Family{
				Codes:          []uint64{0x1bd8a64ad10},
				NBits:          41,
				Name:           "Standard41h12",
				Hamming:        12,
				TotalWidth:     9,
				WidthAtBorder:  5,
				ReversedBorder: true,
			},
		},

		"Standard52h13": {
			NCodes: 48714,
			Family: Family{
				Codes:          []uint64{0x4064a19651ff1},
				NBits:          52,
				Name:           "Standard52h13",
				Hamming:        13,
				TotalWidth:     10,
				WidthAtBorder:  6,
				ReversedBorder: true,
			},
		},
	}

	for name, expected := range testData {
		comment := fmt.Sprintf("testing %s ", name)
		family, err := GetFamily(name)
		if s.NoError(err, comment+" error") == false {
			continue
		}
		s.Equal(expected.Family.Name, family.Name, comment+" name")
		// Not using .Len as the array are all big and meaningless
		if s.Equal(expected.NCodes, len(family.Codes), comment+" codes") == true {
			s.Equal(family.Codes[0], expected.Family.Codes[0], comment+" codes[0]")
		}
		s.Equal(expected.Family.NBits, family.NBits, comment+" codes")
		s.Equal(expected.Family.Hamming, family.Hamming, comment+" hamming")
		s.Equal(expected.Family.TotalWidth, family.TotalWidth, comment+" TotalWidth")
		s.Equal(expected.Family.WidthAtBorder, family.WidthAtBorder, comment+" WidthAtBorder")
		s.Equal(expected.Family.ReversedBorder, family.ReversedBorder, comment+" ReversedBorder")
	}

}

func requirePngGrayImage(s *suite.Suite, path string) *image.Gray {
	comment := fmt.Sprintf("Reading '%s'", path)
	file, err := os.Open(path)
	s.Require().NoError(err, comment)
	defer func() { s.Require().NoError(file.Close()) }()
	img, err := png.Decode(file)
	s.Require().NoError(err, comment)
	res := image.NewGray(img.Bounds())
	draw.Draw(res, res.Bounds(), img, image.Pt(0, 0), draw.Src)
	return res
}

func (s *FamilySuite) TestTagRendering() {
	officialFamilies := []string{"16h5", "25h9", "36h11", "Circle21h7",
		"Circle49h12", "Standard41h12", "Standard52h13"}
	// Note: not testing custom as the transparency causes a lot of
	// issues and we likely do not need it.
	for _, name := range officialFamilies {
		comment := fmt.Sprintf("rendering '%s':0", name)
		expectedImagePath := filepath.Join("utest-data", fmt.Sprintf("tag%s_00000.png", name))
		expected := requirePngGrayImage(&(s.Suite), expectedImagePath)
		family, err := GetFamily(name)
		s.Require().NoError(err, comment)
		s.Equal(expected, family.RenderTag(0), comment)
	}

}

func TestFamilySuite(t *testing.T) {
	suite.Run(t, new(FamilySuite))
}
