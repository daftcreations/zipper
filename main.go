package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"code.cloudfoundry.org/bytefmt"
)

var wg sync.WaitGroup

type filess struct {
	name string
	size int
}

func main() {
	var (
		err               error
		newFiless         []filess
		compressFilesList []string
		tmpZipSplitSize   int
		dirPath           string
	)
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Set path postfix as per OS
	pathPostFix := `/`
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
		tmpZipSplitSize = 5000 // 5MB
	}
	// Change to Bytes
	zipSplitSize := (tmpZipSplitSize - 400) * 1000
	log.Println("Splitting into", tmpZipSplitSize*1000)

	// For folder/directory path
	if len(os.Args) >= 3 {
		dirPath = os.Args[2]
	} else {
		fmt.Println("Making zip of " +
			bytefmt.ByteSize(uint64(tmpZipSplitSize*1000))) // returns "1K")
		fmt.Println("---------------------")

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter path to folder/directory or drag and drop that folder/directory here")
		fmt.Println("---------------------")

		text, _ := reader.ReadString('\n')
		dirPath = strings.Replace(text, "\n", "", -1)
	}
	dir, err := os.Stat(dirPath)
	if err != nil {
		log.Println("failed to open directory, error:", err)
	}
	if !dir.IsDir() {
		fmt.Println("Kindly enter absolute path of folder/directory")
		log.Fatalf("%q is not a directory", dir.Name())
	}
	err = filepath.Walk(dirPath,
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
					int(info.Size()),
				})
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	var fileSizeSum int = 0
	var pass int = 0

	// dirSize, _ := DirSize(dirPath)
	for _, v := range newFiless {
		if fileSizeSum <= zipSplitSize {
			compressFilesList = append(compressFilesList, v.name)
			fileSizeSum += v.size
		} else {
			log.Println("Pass", pass)
			tmpCompress, err := ioutil.TempDir("", "prefix")
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range compressFilesList {
				log.Println("Moving", v, "to", tmpCompress)
				copy(
					v, // source
					filepath.Join(tmpCompress, filepath.Base(v)), // dest tmp/<filename>
					int64(fileSizeSum/len(compressFilesList)))    // buffer
			}
			if err := zipSource(tmpCompress+pathPostFix,
				path.Dir(dirPath)+pathPostFix+dir.Name()+`-`+fmt.Sprint(pass+1)+".zip"); err != nil {
				log.Println("Error compressing", tmpCompress, err)
			} else {
				log.Println("Compressing", tmpCompress)
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
func copy(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s already exists.", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

// Size of dir
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
