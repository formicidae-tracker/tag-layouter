package main

import (
	"gihtub.com/formicidae-tracker/tag-layouter/internal/tag"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Family tag.FamilyBlock `short:"t" long:"family-and-size" description:"Family and size to use. format: 'name:size:begin-end'" required:"yes"`
	Number int             `short:"n" long:"number" description:"Number of tags to display in an arena" required:"yes" `
	Width  float64         `short:"W" long:"width" description:"Width to use in mm" default:"210"`
	Height float64         `short:"H" long:"height" description:"Height to use in mm" default:"297"`
	Margin float64         `short:"m" long:"margin" description:"For paper in mm " default:"20.0"`
	DPI    int             `short:"d" long:"dpi" description:"DPI to use" default:"300"`

	Args struct {
		File flags.Filename
	} `positional-args:"yes" required:"yes"`
}
