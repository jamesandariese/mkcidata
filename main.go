package main

import (
	"log"
	"os"

	mkiso "github.com/jamesandariese/mkcidata/pkg"

	"flag"
)

var (
	out string
)

func main() {
	flag.Parse()
	if flag.NArg() < 3 {
		log.Printf("USAGE: %s output.iso user-data meta-data", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}
	err := mkiso.CreateIso(flag.Args()[0], flag.Args()[1:])
	if err != nil {
		log.Fatalf("Fatal error making ISO: %v", err)
	}
}
