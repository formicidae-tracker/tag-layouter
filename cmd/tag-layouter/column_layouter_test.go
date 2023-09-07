package main

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/formicidae-tracker/tag-layouter/internal/tag"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ColumnLayouterSuite struct {
	suite.Suite
	drawer *MockVectorDrawer
}

func (s *ColumnLayouterSuite) SetupTest() {
	s.drawer = NewMockVectorDrawer(s.T())
}

func (s *ColumnLayouterSuite) TearDownTest() {
	s.drawer = nil
}

func TestColumnLayouterSuite(t *testing.T) {
	fmt.Printf("test suite\n")
	suite.Run(t, new(ColumnLayouterSuite))
}

func mustFamily(f *tag.Family, err error) *tag.Family {
	if err != nil {
		panic(err.Error())
	}
	if f == nil {
		panic("nil family")
	}
	return f
}

func (s *ColumnLayouterSuite) TestSingleRow() {
	block := PlacedBlock{
		FamilyBlock: tag.FamilyBlock{
			Family: mustFamily(tag.GetFamily("Standard41h12")),
			SizeMM: 1.6,
			Ranges: []tag.Range{{Begin: 1, End: 11}},
		},
		Height:            400,
		Width:             400,
		X:                 0,
		Y:                 0,
		ActualTagWidth:    18,
		ActualBorderWidth: 2,
		CutLineWidth:      0,
		NTagsPerRow:       10,
		Skips:             0,
		DPI:               300,
	}

	var calls []*mock.Call
	for i := 0; i < 10; i++ {
		translate := s.drawer.EXPECT().TranslateScale(image.Pt(2+i*20, 2), 2).Return().Once()
		draw := s.drawer.EXPECT().DrawPolygons(mock.Anything).Return().NotBefore(translate).Once()
		end := s.drawer.EXPECT().EndTranslate().Return().NotBefore(draw).Once()
		calls = append(calls, translate, draw, end)
	}
	s.drawer.EXPECT().Label(image.Pt(2, 4), "foo", mock.Anything, color.Gray{}).Return().Once().NotBefore(calls...)
	s.NoError(block.Render(s.drawer, "foo"))

}
