package main

import (
	"errors"
	"log"
)

func main() {
	if err := execute(); err != nil {
		log.Fatalf("unhandled error: %s", err)
	}
}

func execute() error {
	return errors.New("Not Yet Implemented")
}
