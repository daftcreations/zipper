package zipper

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/InVisionApp/tabular"
	"github.com/mholt/archiver"
	"github.com/oleiade/lane"
	. "github.com/pratikbalar/zipper/pkg"
)

var (
	wg           sync.WaitGroup
	tab          tabular.Table
	osPathSuffix string = "/"
)

type zipTask struct {
	zipFileList []string
	dest        string
	count       int
	format      string
}

var (
	count        int = 1
	filesList    []string
	goidRawSize  int = 10
	filesRowSize int = 25
	zipRowSize   int = 25
)

func init() {
	tab = tabular.New()
	tab.Col("goid", "Workerid", goidRawSize)
	tab.Col("files", "Files", filesRowSize)
	tab.Col("zip", "Zip", zipRowSize)

	if runtime.GOOS == `windows` {
		osPathSuffix = `\`
	}
}

// Create zip size of zipSplitSize bytes from give dirPath Directory
func CrateZips(dirPath string, zipSplitSize int) error {
	// Check if dirPath is empty
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		log.Fatal("Directory is empty")
	}

	now := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())
	dirPath = strings.TrimRight(dirPath, osPathSuffix)
	zipQueue := make(chan zipTask, runtime.NumCPU()*4)
	noOfWorker := runtime.NumCPU()

	// Spin up #noOfWorker workers
	wg.Add(noOfWorker)
	for i := 0; i < noOfWorker; i++ {
		go makeArchive(zipQueue)
	}

	queue := lane.NewQueue()

	if err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Ignore Directory
			if info.IsDir() {
				return nil
			}

			// If file size is less than zipSplitSize, then throw error
			if info.Size() > int64(zipSplitSize) {
				return fmt.Errorf("\"%v\" is \"%vKB\" more then %vKB",
					info.Name(), info.Size()/1000, zipSplitSize/1000)
			}

			// Add file in queue to be processed by workers
			queue.Enqueue(filepath.Join(dirPath, info.Name()))

			return nil
		}); err != nil {
		return fmt.Errorf("error walking through path: %v", err)
	}

	totalBytes := 0
	buf := *new(bytes.Buffer)
	zipWriter := zip.NewWriter(&buf)
	format := tab.Print("*")

	for {
		singleFile := fmt.Sprint(queue.Dequeue())
		filesList = append(filesList, singleFile)

		// EOF Predict zip size
		zipFile, err := zipWriter.Create(filepath.Base(singleFile))
		if err != nil {
			return err
		}
		fileBody, err := os.ReadFile(filepath.Clean(singleFile))
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
		// EOF Predict zip size

		totalBytes += zippedFileSize

		if totalBytes > zipSplitSize || queue.Empty() {
			// Add task to queue if its not empty
			if !queue.Empty() {
				queue.Enqueue(singleFile)
			}

			// Add last task to queue channel, close the channel and break the loop
			if queue.Empty() {
				zipQueue <- zipTask{
					filesList,
					fmt.Sprintf("%s-%v.zip", filepath.Base(dirPath), count),
					count,
					format,
				}
				close(zipQueue)
				break
			}

			// Add task to zip queue channel
			zipQueue <- zipTask{
				filesList[:len(filesList)-1],
				fmt.Sprintf("%s-%v.zip", filepath.Base(dirPath), count),
				count,
				format,
			}

			filesList = []string{}
			totalBytes = 0
			count++
		}
	}
	// If channel is
	if _, ok := <-zipQueue; ok {
		close(zipQueue)
	}

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("Fin. Took ", time.Since(now))
	return nil
}

func makeArchive(zipQueue chan zipTask) {
	for {
		zipTask, ok := <-zipQueue
		if !ok {
			wg.Done()
			return
		}
		if err := archiver.Archive(zipTask.zipFileList, zipTask.dest); err != nil {
			log.Fatal(err)
		}
		fmt.Printf(zipTask.format, Goid(), filepath.Base(zipTask.zipFileList[0]), zipTask.dest)
		for _, v := range zipTask.zipFileList[1:] {
			fmt.Printf(zipTask.format, "", filepath.Base(v), "")
		}
		// Adding gap in first row of table
		for i := 0; i < zipRowSize+filesRowSize+goidRawSize+2; i++ {
			fmt.Print("-")
		}
		fmt.Println()
	}
}
