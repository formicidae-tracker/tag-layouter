package main

type Options struct {
	FamilyAndSize string  `short:"t" long:"family-and-size" description:"Families and size to use. format: 'name:size:begin-end'"`
	Number        int     `short:"n" long:"number" description:"Number of tags to display in an arena" default:"0"`
	Width         float64 `short:"W" long:"width" description:"Width to use in mm" default:"210"`
	Height        float64 `short:"H" long:"height" description:"Height to use in mm" default:"297"`
	Margin        float64 `short:"m" long:"margin" description:"For paper in mm " default:"20.0"`
	DPI           int     `short:"d" long:"dpi" description:"DPI to use" default:"300"`
}
