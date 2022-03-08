package zipper

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/InVisionApp/tabular"
	"github.com/mholt/archiver"
	"github.com/oleiade/lane"
	. "github.com/pratikbalar/zipper/pkg"
)

var (
	wg  sync.WaitGroup
	tab tabular.Table
)

type zipTask struct {
	zipFileList []string
	dest        string
	count       int
	format      string
}

type filess struct {
	name string
	size int64
}

var (
	newFiless    []filess
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
}

func CrateZips(dirPath string, zipSplitSize int) error {
	zipQueue := make(chan zipTask, runtime.NumCPU()*4)
	noOfWorker := runtime.NumCPU()
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
			if info.IsDir() {
				return nil
			}

			if info.Size() > int64(zipSplitSize) {
				return fmt.Errorf("\"%v\" is \"%vKB\" more then %vKB",
					info.Name(), info.Size()/1000, zipSplitSize/1000)
			}

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
	if _, ok := <-zipQueue; ok {
		close(zipQueue)
	}
	wg.Wait()
	fmt.Println("Fin.")
	return nil
}

func makeArchive(zipQueue chan zipTask) {
	for {
		zipTask, ok := <-zipQueue
		if ok == false {
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
		for i := 0; i < zipRowSize+filesRowSize+goidRawSize+2; i++ {
			fmt.Print("-")
		}
		fmt.Println()
	}
}
