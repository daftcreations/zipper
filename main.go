package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

func main() {
	pathPostFix := `\`
	if runtime.GOOS == `windows` {
		pathPostFix = `/`
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(pwd)
	thisRegex := regexp.MustCompile(`^this [0-9]+`)

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	var newFiless []filess
	var compressFilesList []string
	for _, f := range files {
		if f.IsDir() && thisRegex.MatchString(f.Name()) {

			// Accept in KB
			tmpZipSplitSize, err := strconv.Atoi(strings.Split(f.Name(), " ")[1])
			if err != nil {
				log.Println(err)
			}

			// Change to Bytes
			zipSplitSize := (tmpZipSplitSize - 500) * 1000
			err = filepath.Walk(f.Name(),
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					// Skip dir
					if info.IsDir() {
						return nil
					}
					newFiless = append(
						newFiless,
						filess{
							filepath.Join(f.Name(), info.Name()),
							int(info.Size()),
						})
					return nil
				})
			if err != nil {
				log.Println(err)
			}
			var fileSizeSum int = 0
			var i int = 0
			for _, v := range newFiless {
				if fileSizeSum <= zipSplitSize {
					compressFilesList = append(compressFilesList, v.name)
					fileSizeSum += v.size
				} else {
					log.Println("Pass", i)
					tmpCompress, err := ioutil.TempDir("", "prefix")
					if err != nil {
						log.Fatal(err)
					}
					for _, v := range compressFilesList {
						log.Println("Moving", v, "to", tmpCompress)
						wg.Add(1)
						go copy(
							v, // source
							filepath.Join(tmpCompress, filepath.Base(v)), // dest tmp/<filename>
							int64(fileSizeSum/len(compressFilesList)))    // buffer
						wg.Wait()
					}
					wg.Add(1)
					go zipSource(tmpCompress+pathPostFix, fmt.Sprint(i)+".zip")
					wg.Wait()
					fileSizeSum = 0
					compressFilesList = nil
					i++
				}
			}
			break
		}
	}
}

func zipSource(source, target string) error {
	defer wg.Done()
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

func copy(src, dst string, BUFFERSIZE int64) error {
	defer wg.Done()
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
