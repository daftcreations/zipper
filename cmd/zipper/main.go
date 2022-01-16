package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	. "github.com/mholt/archiver/v3"
)

var wg sync.WaitGroup

type filess struct {
	name string
}

var (
	newFiless         []filess
	compressFilesList []string
	pathPostFix       string = "/"
	pwd               string
)

func main() {
	var (
		err             error
		tmpZipSplitSize int
	)

	// This utility will mostly run on low power CPUs so
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Set path postfix as per OS
	if runtime.GOOS == `windows` {
		pathPostFix = `\`
	}
	// For zip size
	if len(os.Args) >= 2 {
		if tmpZipSplitSize, err = strconv.Atoi(os.Args[1]); err != nil {
			log.Fatalln("Converting tmpZipSplitSize: ", err)
		}
	} else {
		tmpZipSplitSize = 3000 // 5MB
	}
	// Change to Bytes
	zipSplitSize := tmpZipSplitSize * 1000
	log.Println("Splitting into", zipSplitSize)

	// Trim ending linux and mac `/` or in windows `\` from path
	if err := crateZips(
		strings.TrimRight(os.Args[2], pathPostFix), zipSplitSize); err != nil {
		log.Fatal("Error creating zip: ", err)
	}
}

func crateZips(dirPath string, zipSplitSize int) error {
	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Skip dir
			if info.IsDir() {
				return nil
			}

			// Exit if any file is more then zipSplitSize
			if info.Size() > int64(zipSplitSize) {
				return fmt.Errorf("\"%v\" is more then %vKB\n", info.Name(), zipSplitSize)
			}
			return nil
		})
	if err != nil {
		return fmt.Errorf("Error walking through path: %v", err)
	}

	_ = filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Skip dir
			if info.IsDir() {
				return nil
			}

			// Append files in struct list
			newFiless = append(
				newFiless,
				filess{
					filepath.Join(dirPath, info.Name()),
				})
			return nil
		})

	// Adding last file seondtime to handle last file case when creating zip
	for k, v := range newFiless {
		if k == len(newFiless)-1 {
			newFiless = append(newFiless, v)
		}
	}

	var (
		count     int = 1
		filesList []string
	)
	// ch := make(chan struct{}, runtime.NumCPU())

	for k, v := range newFiless {
		filesList = append(filesList, v.name)
		tmpCompress, err := ioutil.TempDir("", "prefix")
		if err != nil {
			return fmt.Errorf("Error creating temp dir: %v", err)
		}
		tmpZipPath := filepath.Join(tmpCompress, "test.zip")
		if err = Archive(filesList, tmpZipPath); err != nil {
			return fmt.Errorf("Error creating temp archive: %v", err)
		}
		fileStat, err := os.Stat(tmpZipPath)
		if err != nil {
			return fmt.Errorf("Error getting file info of tmpzip file: %v", err)
		}
		if err = os.RemoveAll(tmpZipPath); err != nil {
			return fmt.Errorf("Error removing temp zipfile: %v", err)
		}

		if fileStat.Size() > int64(zipSplitSize) || k == len(newFiless)-1 {
			wg.Wait()
			wg.Add(1)
			go func(filesList []string, dest string) {
				fmt.Println("PASS", fmt.Sprint(count), ": Creating",
					filepath.Join(dirPath, dest), ": -----------------------------")
				if err := Archive(filesList, dest); err != nil {
					log.Fatal(err)
				}
				wg.Done()
			}(filesList[:len(filesList)-1], filepath.Base(dirPath)+"-"+fmt.Sprint(count)+".zip")

			// log.Println("Cretain zip of ", filesList)
			lastFile := filesList[len(filesList)-1]
			filesList = []string{}
			filesList = append(filesList, lastFile)

			count++
		}
		// log.Println("filesList", filesList)
		wg.Wait()
	}
	return nil
}
