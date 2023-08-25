package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	if filepath.Ext(string(opts.Args.File)) != ".svg" {
		return fmt.Errorf("invalid filepath '%s': only SVG are supported, filepath must end with '.svg'", opts.Args.File)
	}

	f, err := os.Create(string(opts.Args.File))
	if err != nil {
		return err
	}
	defer f.Close()
	SVG := svg.New(f)

	return opts.LayoutArena(SVG)
}
