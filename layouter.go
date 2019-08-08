package main

type Layouter interface {
	Layout(drawer Drawer, families []FamilyBlock) error
}
