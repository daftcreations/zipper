package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/mholt/archiver/v3"
)

var wg sync.WaitGroup

type filess struct {
	name string
	size int
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
		dirPath         string
	)
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Set path postfix as per OS
	if runtime.GOOS == `windows` {
		pathPostFix = `\`
	}

	// Get current path of main
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Current working dir:", pwd)

	// For zip size
	if len(os.Args) >= 2 {
		if tmpZipSplitSize, err = strconv.Atoi(os.Args[1]); err != nil {
			log.Fatalln("Converting tmpZipSplitSize: ", err)
		}
	} else {
		tmpZipSplitSize = 3000 // 5MB
	}
	// Change to Bytes
	zipSplitSize := (tmpZipSplitSize - 400) * 1000
	log.Println("Splitting into", tmpZipSplitSize*1000)

	//  Trim ending linux and mac `/` or in windows `\` from path
	dirPath = os.Args[2]
	crateZips(dirPath, zipSplitSize)
}

func crateZips(dirPath string, zipSplitSize int) {
	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Skip dir
			if info.IsDir() {
				return nil
			}

			// Exit if any file is more then tmpZipSplitSize
			if info.Size() > int64(zipSplitSize) {
				fmt.Println("####################")
				fmt.Printf("\"%v\" is more then %vKB\n", info.Name(), zipSplitSize)
				fmt.Println("####################")
				os.Exit(1)
			}

			// Append files in struct list
			newFiless = append(
				newFiless,
				filess{
					filepath.Join(dirPath, info.Name()),
					int(info.Size()),
				})
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	// Disable sorting for now
	// sort.SliceStable(newFiless, func(i, j int) bool {
	// 	return newFiless[i].size < newFiless[j].size
	// })

	// log.Println(newFiless)

	var fileSizeSum int = 0
	var pass int = 0
	// dirSize, _ := DirSize(dirPath)
	for k, v := range newFiless {
		log.Println(k, ":", v)
		compressFilesList = append(compressFilesList, v.name)
		fileSizeSum += v.size
		if fileSizeSum < zipSplitSize {
			log.Println(v.name)
			if k+1 == len(newFiless) {
				log.Println("Pass", pass, fileSizeSum, "<", zipSplitSize)
				tmpCompress, err := ioutil.TempDir("", "prefix")
				if err != nil {
					log.Fatal(err)
				}
				for _, v := range compressFilesList {
					log.Println("Moving", v, "to", tmpCompress)
					err = copy(v, filepath.Join(tmpCompress, filepath.Base(v)))
					if err != nil {
						log.Fatalf("Cannot copy \"%v\" to \"%v\"\n",
							v, filepath.Join(tmpCompress, filepath.Base(v)))
					}
				}
				destZipPath := filepath.Join(
					strings.TrimRight(dirPath, filepath.Base(dirPath)),
					fmt.Sprint(pass+1)+".zip")
				err = archiver.Archive(
					[]string{tmpCompress + pathPostFix},
					destZipPath)
				if err != nil {
					log.Println("Error compressing", tmpCompress, err)
				} else {
					log.Println("Compressing", tmpCompress)
					log.Println("Creating zip at", destZipPath)
				}
				if err := os.RemoveAll(tmpCompress); err != nil {
					log.Println("Error removing", tmpCompress, err)
				} else {
					log.Println("Removing", tmpCompress)
				}
				fileSizeSum = 0
				compressFilesList = nil
				pass++
			}
		} else {
			log.Println("Pass", pass, fileSizeSum, "<", zipSplitSize)
			tmpCompress, err := ioutil.TempDir("", "prefix")
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range compressFilesList {
				log.Println("Moving", v, "to", tmpCompress)
				err = copy(v, filepath.Join(tmpCompress, filepath.Base(v)))
				if err != nil {
					log.Fatalf("Cannot copy \"%v\" to \"%v\"\n",
						v, filepath.Join(tmpCompress, filepath.Base(v)))
				}
			}
			destZipPath := filepath.Join(
				strings.TrimRight(dirPath, filepath.Base(dirPath)),
				fmt.Sprint(pass+1)+".zip")
			err = archiver.Archive(
				[]string{tmpCompress + pathPostFix},
				destZipPath)
			if err != nil {
				log.Println("Error compressing", tmpCompress, err)
			} else {
				log.Println("Compressing", tmpCompress)
				log.Println("Creating zip at", destZipPath)
			}
			if err := os.RemoveAll(tmpCompress); err != nil {
				log.Println("Error removing", tmpCompress, err)
			} else {
				log.Println("Removing", tmpCompress)
			}
			fileSizeSum = 0
			compressFilesList = nil
			pass++
		}
	}
}

// Copy file
func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
