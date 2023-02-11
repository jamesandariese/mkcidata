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
	args := flag.Args()
	if len(args) < 2 {
		log.Fatalf("usage: mkcidata thing.iso filea [fileb] [filec]")
	}
	err := mkiso.CreateIso(args[0], args[1:])
	if err != nil {
		log.Fatalf("Fatal error making ISO: %v", err)
	}
}
