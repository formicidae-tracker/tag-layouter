package main

import "fmt"

type ColumnLayouter struct {
	Width        float64
	Height       float64
	Columns      int
	FamilyMargin float64
	TagBorder    float64
	PaperBorder  float64
}

func (c *ColumnLayouter) Layout(drawer Drawer, families []FamilyAndSize) error {
	return fmt.Errorf("Not yet implemneted")
}
