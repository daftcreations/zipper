package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pkg/profile"
	. "github.com/pratikbalar/zipper/internal/zipper"
)

// CLI flags variables
var (
	profEnable      string = "false"
	tmpZipSplitSize int

	Version   = "0.0.0"
	CommitSha = "xxxx"
	BuildTime = "0000-00-00T00:00:00+00:00"
)

func main() {
	var (
		err  error
		path string

		// CLI flags
		help       bool
		getVersion bool
	)

	flag.BoolVar(&getVersion, "version", false, "get version")
	flag.BoolVar(&help, "help", false, "display help and exit")

	flag.Parse()

	if getVersion {
		fmt.Printf(`zipper v%s (sha: %s) (BuildTime: %s)

Made by DaftCreations, Opensource org For People/Creator, By People/Creator.

Soon to be DAO - https://github.com/daftcreations
`,
			Version, CommitSha, BuildTime)
		return
	}

	if help || len(flag.Args()) == 0 {
		fmt.Print(`zipper Usage:

$ zipper <size> # size in KiloBytes
Splitting into <size> KB

Enter path you want to zip:
> <PATH>

Flags:
`)
		flag.PrintDefaults()
		return
	}

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
		fmt.Print("> ")
		if _, err = fmt.Scanln(&path); err != nil {
			log.Fatal(err)
		}
		err := CrateZips(path, zipSplitSize)
		if err != nil {
			log.Fatal("Error creating zip: ", err)
		}
	}
}
