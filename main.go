package main

import (
	"log"

	mkiso "github.com/jamesandariese/mkcidata/pkg"

	"flag"
)

var (
	out string
)

func main() {
	flag.Parse()
	err := mkiso.CreateIso("poop.iso", flag.Args())
	if err != nil {
		log.Fatalf("Fatal error making ISO: %v", err)
	}
}
