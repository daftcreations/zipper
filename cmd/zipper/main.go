package main

import (
	"archive/zip"
	"crypto/sha1"
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

func init() {
	h := sha1.New()
	h.Write([]byte(os.Getenv("UNLOCK_ZIPPER")))

	if fmt.Sprintf("%x", h.Sum(nil)) != "f43dc14eb92b5dafba73481449f8ed3a602c0b6e" {
		os.Exit(2)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Set path postfix as per OS
	if runtime.GOOS == `windows` {
		pathPostFix = `\`
	}
}

func main() {
	var (
		err             error
		tmpZipSplitSize int
		// dirPath         string
	)

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

	// Trim ending linux and mac `/` or in windows `\` from path
	crateZips(strings.TrimRight(os.Args[2], pathPostFix), zipSplitSize)
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
				moveAndZip(&pass, &fileSizeSum, zipSplitSize, dirPath)
			}
		} else {
			moveAndZip(&pass, &fileSizeSum, zipSplitSize, dirPath)
		}
	}
}

// MoveandZip
func moveAndZip(pass, fileSizeSum *int, zipSplitSize int, dirPath string) {
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
		filepath.Base(dirPath)+"-"+fmt.Sprint(*pass+1)+".zip")
	if err = zipSource(tmpCompress+pathPostFix, destZipPath); err != nil {
		log.Println("Error compressing", tmpCompress, err)
	}
	log.Println("Compressing", tmpCompress)
	log.Println("Creating zip at", destZipPath)
	if err := os.RemoveAll(tmpCompress); err != nil {
		log.Println("Error removing", tmpCompress, err)
	}
	log.Println("Removing", tmpCompress)
	*fileSizeSum = 0
	compressFilesList = nil
	*pass++
}

// Create zip
func zipSource(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
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
