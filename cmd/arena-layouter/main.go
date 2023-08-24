package main

import (
	"log"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/jessevdk/go-flags"
)

func main() {
	if err := execute(); err != nil {
		log.Fatalf("unhandled error: %s", err)
	}
}

func execute() error {
	var opts Options
	if _, err := flags.Parse(&opts); err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if err := opts.Validate(); err != nil {
		return err
	}

	f, err := os.Create(string(opts.Args.File))
	if err != nil {
		return err
	}
	defer f.Close()
	SVG := svg.New(f)

	opts.LayoutArena(SVG)

	return nil
}
