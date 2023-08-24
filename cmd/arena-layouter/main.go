package main

import (
	"errors"
	"fmt"
	"log"
	"os"

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
	fmt.Printf("%+v", opts)
	return errors.New("not yet implemented")
}
