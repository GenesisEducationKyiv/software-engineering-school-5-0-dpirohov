package config

import (
	"log"
	"os"
)

var ROOT_DIR string

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get working directory: %v", err)
	}

	ROOT_DIR = dir
}
