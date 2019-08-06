package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"sort"
)

type ColumnLayouter struct {
	Width        float64
	Height       float64
	Columns      int
	FamilyMargin float64
	TagBorder    float64
	PaperBorder  float64
	drawer       Drawer
}

func (c *ColumnLayouter) PerfectPixelSizeMM(s float64, b float64, totalWidth int) (tagSize float64, borderSize float64) {
	perfectPixelSize := c.drawer.ToDot(s / float64(totalWidth))
	tagSize = c.drawer.ToMM(perfectPixelSize) * float64(totalWidth)
	borderSize = c.drawer.ToMM(int(math.Round(float64(totalWidth)*b)) * perfectPixelSize)
	return tagSize, borderSize
}

func (f *FamilyAndSize) FamilyLabel() string {
	return fmt.Sprintf("%s %.2fMM", f.Family.Name, f.Size)
}

func (c *ColumnLayouter) FamilyHeight(f FamilyAndSize, columnWidth float64) float64 {
	actualTagWidth, actualBorderWidth := c.PerfectPixelSizeMM(f.Size, c.TagBorder, f.Family.TotalWidth)
	nbSlots := f.End - f.Begin - 1 + len(f.FamilyLabel())*2/3

	nbTagsPerRow := int(math.Floor(columnWidth / (actualTagWidth + actualBorderWidth)))
	nbRows := nbSlots / nbTagsPerRow
	if nbSlots%nbTagsPerRow != 0 {
		nbRows += 1
	}
	return float64(nbRows)*(actualTagWidth+actualBorderWidth) - actualBorderWidth
}

type FamilyWithHeight struct {
	FamilyAndSize
	Height float64
}

type FamilyWithHeightList []FamilyWithHeight

func (fhs FamilyWithHeightList) Len() int {
	return len(fhs)
}

func (fhs FamilyWithHeightList) Less(i, j int) bool {
	return fhs[i].Height < fhs[j].Height
}

func (fhs FamilyWithHeightList) Swap(i, j int) {
	fhs[i], fhs[j] = fhs[j], fhs[i]
}

func (c *ColumnLayouter) LayoutOne(xOffset, yOffset float64, f FamilyWithHeight, columnWidth float64) {

	label := f.FamilyLabel()
	actualTagWidth, actualBorderWidth := c.PerfectPixelSizeMM(f.Size, c.TagBorder, f.Family.TotalWidth)
	log.Printf("%s:%.2fmm actual size: %.2f; error: %.2f%%",
		f.Family.Name,
		f.Size,
		actualTagWidth,
		math.Abs(actualTagWidth-f.Size)/f.Size*100)
	nbTagsPerRow := int(math.Floor(columnWidth / (actualTagWidth + actualBorderWidth)))

	skips := len(label) * 2 / 3

	ix := skips % nbTagsPerRow
	iy := skips / nbTagsPerRow
	for i := f.Begin; i < f.End; i++ {
		x := float64(ix)*(actualTagWidth+actualBorderWidth) + xOffset
		y := float64(iy)*(actualTagWidth+actualBorderWidth) + yOffset
		DrawTag(c.drawer, f.Family, f.Family.Codes[i], x, y, f.Size, 0, nil)
		ix += 1
		if ix >= nbTagsPerRow {
			ix = 0
			iy += 1
		}
	}
	c.drawer.Label(c.drawer.ToDot(xOffset), c.drawer.ToDot(yOffset+actualTagWidth-actualBorderWidth/2), actualTagWidth, label, color.RGBA{0xff, 00, 00, 0xff})

}

func (c *ColumnLayouter) Layout(drawer Drawer, families []FamilyAndSize) error {
	c.drawer = drawer
	if c.Columns < 1 {
		return fmt.Errorf("Invalid number of column")
	}

	c.FamilyMargin = drawer.ToMM(drawer.ToDot(c.FamilyMargin))
	c.PaperBorder = drawer.ToMM(drawer.ToDot(c.PaperBorder))

	c.drawer.DrawRectangle(0, 0, c.drawer.ToDot(c.Width), c.drawer.ToDot(c.Height), color.White)

	columnWidth := (c.Width - 2*c.PaperBorder - c.FamilyMargin*float64(c.Columns-1)) / float64(c.Columns)
	columnWidth = drawer.ToMM(drawer.ToDot(columnWidth))
	columnHeight := c.Height - 2*c.PaperBorder

	cFamilies := []FamilyWithHeight{}

	for _, f := range families {
		h := c.FamilyHeight(f, columnWidth)
		cFamilies = append(cFamilies, FamilyWithHeight{
			FamilyAndSize: f,
			Height:        h,
		})
	}
	sort.Sort(sort.Reverse(FamilyWithHeightList(cFamilies)))
	columns := make([][]FamilyWithHeight, c.Columns)

	for _, fh := range cFamilies {
		fitted := false
		for idxCol, col := range columns {
			height := 0.0
			for _, fhc := range col {
				height += fhc.Height + c.FamilyMargin
			}
			if fh.Height+height > columnHeight {
				continue
			}
			columns[idxCol] = append(columns[idxCol], fh)
			fitted = true
			break
		}
		if fitted == false {
			return fmt.Errorf("Could not fill %s:%.2f:%d-%d in layout", fh.Family.Name, fh.Size, fh.Begin, fh.End)
		}
	}

	x := c.PaperBorder
	for _, column := range columns {
		y := c.PaperBorder
		for _, f := range column {
			c.LayoutOne(x, y, f, columnWidth)
			y += f.Height + c.FamilyMargin
		}
		x += columnWidth + c.FamilyMargin
	}

	return nil
}
