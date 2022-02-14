package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/profile"
	. "github.com/pratikbalar/zipper/internal/zipper"
)

var (
	profEnable      string = "false"
	osPathSuffix    string = "/"
	tmpZipSplitSize int
)

func main() {
	var err error

	if profEnable == "true" {
		defer profile.Start(profile.ProfilePath("."),
			profile.MemProfile, profile.MemProfileRate(1),
			// profile.CPUProfile,
			// profile.TraceProfile,
		).Stop()
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if runtime.GOOS == `windows` {
		osPathSuffix = `\`
	}
	if len(os.Args) >= 2 {
		if tmpZipSplitSize, err = strconv.Atoi(os.Args[1]); err != nil {
			log.Fatalln("Converting tmpZipSplitSize: ", err)
		}
	}
	zipSplitSize := tmpZipSplitSize * 1000
	fmt.Println("Splitting into", zipSplitSize)

	if err := CrateZips(
		strings.TrimRight(os.Args[2], osPathSuffix), zipSplitSize); err != nil {
		log.Fatal("Error creating zip: ", err)
	}
}
