package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	archiver "github.com/mholt/archiver/v3"
	"github.com/oleiade/lane"
	"github.com/pkg/profile"
)

var wg sync.WaitGroup

type filess struct {
	name string
	size int64
}

type zipTask struct {
	zipFileList []string
	dest        string
	count       int
}

var (
	newFiless       []filess
	osPathSuffix    string = "/"
	profEnable      string = "false"
	count           int    = 1
	filesList       []string
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

	if err := crateZips(
		strings.TrimRight(os.Args[2], osPathSuffix), zipSplitSize); err != nil {
		log.Fatal("Error creating zip: ", err)
	}
}

func crateZips(dirPath string, zipSplitSize int) error {
	zipQueue := make(chan zipTask, runtime.NumCPU()*4)
	noOfWorker := 2
	wg.Add(noOfWorker)
	for i := 0; i < noOfWorker; i++ {
		go makeArchive(zipQueue)
	}

	queue := lane.NewQueue()

	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			if info.Size() > int64(zipSplitSize) {
				return fmt.Errorf("\"%v\" is more then %vKB", info.Name(), zipSplitSize/1000)
			}

			newFiless = append(
				newFiless,
				filess{
					filepath.Join(dirPath, info.Name()), info.Size(),
				})

			return nil
		})
	if err != nil {
		return fmt.Errorf("error walking through path: %v", err)
	}

	sort.Slice(newFiless, func(i, j int) bool {
		return newFiless[i].size < newFiless[j].size
	})

	for _, v := range newFiless {
		queue.Enqueue(v.name)
	}

	totalBytes := 0
	buf := *new(bytes.Buffer)
	zipWriter := zip.NewWriter(&buf)
	for {
		singleFile := fmt.Sprint(queue.Dequeue())
		filesList = append(filesList, singleFile)
		// fmt.Println("File: ", singleFile)

		zipFile, err := zipWriter.Create(filepath.Base(singleFile))
		if err != nil {
			return err
		}
		fileBody, err := os.ReadFile(singleFile)
		if err != nil {
			return err
		}
		zippedFileSize, err := zipFile.Write(fileBody)
		if err != nil {
			return err
		}
		buf.Reset()
		if err = zipWriter.Flush(); err != nil {
			return err
		}

		totalBytes += zippedFileSize

		if totalBytes > zipSplitSize || queue.Empty() {
			if !queue.Empty() {
				queue.Enqueue(singleFile)
			}
			if queue.Size() == 0 {
				zipQueue <- zipTask{
					filesList,
					fmt.Sprintf("%s-%v.zip", filepath.Base(dirPath), count),
					count,
				}
				close(zipQueue)
				break
			}

			zipQueue <- zipTask{
				filesList[:len(filesList)-1],
				fmt.Sprintf("%s-%v.zip", filepath.Base(dirPath), count),
				count,
			}

			filesList = []string{}
			totalBytes = 0
			count++
		}
	}
	if _, ok := <-zipQueue; ok {
		close(zipQueue)
	}
	wg.Wait()
	return nil
}

func makeArchive(zipQueue chan zipTask) {
	for {
		zipTask, ok := <-zipQueue
		if ok == false {
			wg.Done()
			return
		}
		fmt.Printf("%v consuming %v\n", goid(), zipTask)
		// fmt.Println("PASS", fmt.Sprint(count), ": Creating", zipTask.dest, ": -----------------------------")
		// fmt.Println("Creting archive")
		if err := archiver.Archive(zipTask.zipFileList, zipTask.dest); err != nil {
			log.Fatal(err)
		}
	}
}

func goid() int {
	var buf [64]byte
	id, _ := strconv.Atoi(strings.Fields(strings.TrimPrefix(string(buf[:runtime.Stack(buf[:], false)]), "goroutine "))[0])
	return id
}
