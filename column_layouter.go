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
	NColumns     int
	FamilyMargin float64
	TagBorder    float64
	PaperBorder  float64
	CutLine      float64
	drawer       Drawer
}

func (c *ColumnLayouter) PerfectPixelSizeMM(size float64, border float64, cutline float64, totalWidth int) (tagSizeDot int, borderSizeDot int, cutLineSizeDot int) {
	perfectPixelSize := c.drawer.ToDot(size / float64(totalWidth))
	tagSizeDot = perfectPixelSize * totalWidth
	borderSizeDot = int(math.Round(float64(totalWidth)*border)) * perfectPixelSize

	cutLineSizeDot = int(math.Round(float64(borderSizeDot) * cutline))
	if cutline != 0.0 {
		if cutLineSizeDot == 0 {
			cutLineSizeDot = 1
		}
		if (borderSizeDot-cutLineSizeDot)%2 == 1 {
			borderSizeDot += 1
		}
	}
	//	log.Printf("pixel: %d; tag: %d, border: %d, cut %d", perfectPixelSize, perfectPixelSize*totalWidth, borderSizeDot, cutLineSizeDot)

	return tagSizeDot, borderSizeDot, cutLineSizeDot
}

type PlacedFamily struct {
	FamilyBlock
	Height int
	Width  int
	X      int
	Y      int

	ActualTagWidth    int
	ActualBorderWidth int
	CutLineWidth      int
	NTagsPerRow       int
	Skips             int
}

func (c *ColumnLayouter) ComputeFamilySize(f FamilyBlock, columnWidthDot int) PlacedFamily {
	actualTagWidth, actualBorderWidth, cutlineWidth := c.PerfectPixelSizeMM(f.Size, c.TagBorder, c.CutLine, f.Family.TotalWidth)
	skips := len(f.FamilyLabel())*1/2 + 1
	nbSlots := f.NumberOfTags() + skips

	nbTagsPerRow := (columnWidthDot - actualBorderWidth) / (actualTagWidth + actualBorderWidth)
	nbRows := nbSlots / nbTagsPerRow
	if nbSlots%nbTagsPerRow != 0 {
		nbRows += 1
	}

	height := nbRows*(actualTagWidth+actualBorderWidth) + actualBorderWidth
	width := 0
	if nbRows > 1 {
		width = columnWidthDot
	} else {
		width = nbSlots*(actualTagWidth+actualBorderWidth) + actualBorderWidth
	}

	return PlacedFamily{
		FamilyBlock:       f,
		Height:            height,
		Width:             width,
		X:                 0,
		Y:                 0,
		ActualTagWidth:    actualTagWidth,
		ActualBorderWidth: actualBorderWidth,
		CutLineWidth:      cutlineWidth,
		NTagsPerRow:       nbTagsPerRow,
		Skips:             skips,
	}
}

type PlacedFamilyListByHeight []PlacedFamily
type PlacedFamilyListByWidth []PlacedFamily

func (fhs PlacedFamilyListByHeight) Len() int {
	return len(fhs)
}

func (fhs PlacedFamilyListByWidth) Len() int {
	return len(fhs)
}

func (fhs PlacedFamilyListByHeight) Less(i, j int) bool {
	return fhs[i].Height < fhs[j].Height
}

func (fhs PlacedFamilyListByWidth) Less(i, j int) bool {
	return fhs[i].Width < fhs[j].Width
}

func (fhs PlacedFamilyListByHeight) Swap(i, j int) {
	fhs[i], fhs[j] = fhs[j], fhs[i]
}

func (fhs PlacedFamilyListByWidth) Swap(i, j int) {
	fhs[i], fhs[j] = fhs[j], fhs[i]
}

func (c *ColumnLayouter) LayoutOne(pf PlacedFamily) {
	label := pf.FamilyLabel()
	actualSizeMM := c.drawer.ToMM(pf.ActualTagWidth)
	log.Printf("%s:%.2fmm actual size: %.2f; error: %.2f%%",
		pf.Family.Name,
		pf.Size,
		actualSizeMM,
		math.Abs(actualSizeMM-pf.Size)/pf.Size*100)

	ix := pf.Skips % pf.NTagsPerRow
	iy := pf.Skips / pf.NTagsPerRow

	cutLinePos := (pf.ActualBorderWidth - pf.CutLineWidth) / 2
	isFirst := true
	for _, r := range pf.Ranges {
		for i := r.Begin; i < r.End; i++ {
			x := ix*(pf.ActualTagWidth+pf.ActualBorderWidth) + pf.X + pf.ActualBorderWidth
			y := iy*(pf.ActualTagWidth+pf.ActualBorderWidth) + pf.Y + pf.ActualBorderWidth
			DrawTagDot(c.drawer, pf.Family, pf.Family.Codes[i], x, y, pf.ActualTagWidth)
			ix += 1
			if ix >= pf.NTagsPerRow {
				ix = 0
				iy += 1
			}
			if pf.CutLineWidth == 0 {
				continue
			}

			c.drawer.DrawRectangle(x+pf.ActualTagWidth+cutLinePos,
				y,
				pf.CutLineWidth,
				pf.ActualTagWidth,
				color.Black)

			c.drawer.DrawRectangle(x,
				y+pf.ActualTagWidth+cutLinePos,
				pf.ActualTagWidth,
				pf.CutLineWidth,
				color.Black)

			if ix == 1 || isFirst == true {
				isFirst = false
				c.drawer.DrawRectangle(x-pf.CutLineWidth-cutLinePos,
					y,
					pf.CutLineWidth,
					pf.ActualTagWidth,
					color.Black)
			}

			if iy == 0 || iy == 1 && ix <= pf.Skips {
				c.drawer.DrawRectangle(x,
					y-cutLinePos-pf.CutLineWidth,
					pf.ActualTagWidth,
					pf.CutLineWidth,
					color.Black)
			}
		}
	}
	c.drawer.Label(pf.X+pf.ActualBorderWidth, pf.Y, pf.ActualTagWidth, label, color.RGBA{0xff, 00, 00, 0xff})

}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func (c *ColumnLayouter) Layout(drawer Drawer, families []FamilyBlock) error {
	c.drawer = drawer
	if c.NColumns < 1 {
		return fmt.Errorf("Invalid number of column")
	}

	familyMarginDot := drawer.ToDot(c.FamilyMargin)
	paperBorderDot := drawer.ToDot(c.PaperBorder)

	columnWidthDot := (drawer.ToDot(c.Width) - 2*paperBorderDot - familyMarginDot*(c.NColumns-1)) / c.NColumns

	columnHeightDot := drawer.ToDot(c.Height) - 2*paperBorderDot

	placedFamiliesFullWidth := []PlacedFamily{}
	placedFamiliesIncompleteWidth := []PlacedFamily{}
	for _, f := range families {
		pf := c.ComputeFamilySize(f, columnWidthDot)
		if pf.Width < columnWidthDot {
			placedFamiliesIncompleteWidth = append(placedFamiliesIncompleteWidth, pf)
		} else {
			placedFamiliesFullWidth = append(placedFamiliesFullWidth, pf)
		}
	}

	sort.Sort(sort.Reverse(PlacedFamilyListByHeight(placedFamiliesFullWidth)))
	sort.Sort(PlacedFamilyListByWidth(placedFamiliesIncompleteWidth))

	type Column struct {
		Families      []PlacedFamily
		XOffset       int
		Width         int
		Height        int
		LastRowHeight int
	}
	columns := make([]Column, c.NColumns)
	for i, _ := range columns {
		columns[i].XOffset = i*(columnWidthDot+familyMarginDot) + paperBorderDot
		columns[i].Width = 0
		columns[i].Height = 0
		columns[i].LastRowHeight = 0
	}

	for _, pf := range placedFamiliesFullWidth {
		fitted := false
		for idxCol, _ := range columns {
			if (pf.Height + columns[idxCol].Height + familyMarginDot) > columnHeightDot {
				continue
			}
			if len(columns[idxCol].Families) == 0 {
				columns[idxCol].Height = -familyMarginDot
			}
			//			log.Printf("Placing %s in %d %d position", pf.FamilyLabel(), idxCol, len(col.Families))
			pf.X = columns[idxCol].XOffset
			pf.Y = columns[idxCol].Height + familyMarginDot

			columns[idxCol].Families = append(columns[idxCol].Families, pf)
			columns[idxCol].Height += pf.Height + familyMarginDot
			fitted = true
			break
		}
		if fitted == false {
			return fmt.Errorf("Could not fill %s:%.2f:%s in layout", pf.Family.Name, pf.Size, pf.RangeString())
		}
	}

	for _, pf := range placedFamiliesIncompleteWidth {
		fitted := false
		for idxCol, _ := range columns {
			if (pf.Height + columns[idxCol].Height + familyMarginDot) > columnHeightDot {
				//not fitting in height anyway
				continue
			}
			//if we are building a new line
			//check if it fits on the same line
			if (pf.Width + columns[idxCol].Width) > columnWidthDot {
				//no so we terminate the line
				columns[idxCol].Width = 0
				columns[idxCol].Height = columns[idxCol].LastRowHeight
				columns[idxCol].LastRowHeight = 0
				//we recheck if we can be put in height
				if pf.Height+columns[idxCol].Height+familyMarginDot > columnHeightDot {
					continue
				}
			}

			//			log.Printf("Placing small %s in %d %d position", pf.FamilyLabel(), idxCol, len(columns[idxCol].Families))

			pf.X = columns[idxCol].XOffset + columns[idxCol].Width
			pf.Y = columns[idxCol].Height + familyMarginDot

			columns[idxCol].Families = append(columns[idxCol].Families, pf)
			columns[idxCol].Width += pf.Width + familyMarginDot
			columns[idxCol].LastRowHeight = max(columns[idxCol].LastRowHeight, columns[idxCol].Height+familyMarginDot+pf.Height)
			fitted = true
			break
		}
		if fitted == false {
			return fmt.Errorf("Could not fill %s:%.2f:%s in layout", pf.Family.Name, pf.Size, pf.RangeString())
		}
	}

	log.Printf("Filling background")
	c.drawer.DrawRectangle(0, 0, c.drawer.ToDot(c.Width), c.drawer.ToDot(c.Height), color.White)
	log.Printf("Done")

	for _, column := range columns {
		for _, pf := range column.Families {
			c.LayoutOne(pf)
		}
	}

	return nil
}
