package main

import (
	"fmt"
	"testing"

	"github.com/formicidae-tracker/tag-layouter/internal/tag"
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
	fmt.Printf("coucou\n")
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

func (s *ColumnLayouterSuite) Test() {
	block := PlacedBlock{
		FamilyBlock: tag.FamilyBlock{
			Family: mustFamily(tag.GetFamily("Standard41h12")),
			SizeMM: 1.6,
			Ranges: []tag.Range{{Begin: 1, End: -1}},
		},
		Height:         100,
		Width:          100,
		X:              0,
		Y:              0,
		ActualTagWidth: 10,
		CutLineWidth:   0,
		NTagsPerRow:    8,
		Skips:          0,
		DPI:            300,
	}

	block.Render(s.drawer, "")
}
