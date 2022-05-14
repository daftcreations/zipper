package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pkg/profile"
	. "github.com/pratikbalar/zipper/internal/zipper"
)

var (
	profEnable      string = "false"
	tmpZipSplitSize int
)

func main() {
	var err error
	var path string

	if profEnable == "true" {
		defer profile.Start(profile.ProfilePath("."),
			profile.MemProfile, profile.MemProfileRate(1),
			// profile.CPUProfile,
			// profile.TraceProfile,
		).Stop()
	}

	if tmpZipSplitSize, err = strconv.Atoi(os.Args[1]); err != nil {
		log.Fatalln("Converting tmpZipSplitSize: ", err)
	}
	zipSplitSize := tmpZipSplitSize * 1000
	fmt.Println("Splitting into", tmpZipSplitSize, "KB")

	for {
		fmt.Println("\nEnter path you want to zip:")
		fmt.Print(">")
		if _, err = fmt.Scanln(&path); err != nil {
			log.Fatal(err)
		}
		err := CrateZips(path, zipSplitSize)
		if err != nil {
			log.Fatal("Error creating zip: ", err)
		}
	}
}
