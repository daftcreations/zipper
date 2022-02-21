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

	"github.com/mholt/archiver"
	"github.com/oleiade/lane"
	. "github.com/pratikbalar/zipper/pkg"
)

var wg sync.WaitGroup

type zipTask struct {
	zipFileList []string
	dest        string
	count       int
}

type filess struct {
	name string
	size int64
}

var (
	newFiless []filess
	count     int = 1
	filesList []string
)

func CrateZips(dirPath string, zipSplitSize int) error {
	zipQueue := make(chan zipTask, runtime.NumCPU()*4)
	noOfWorker := 2
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

			// newFiless = append(
			// 	newFiless,
			// 	filess{
			// 		filepath.Join(dirPath, info.Name()), info.Size(),
			// 	})

			queue.Enqueue(filepath.Join(dirPath, info.Name()))

			return nil
		}); err != nil {
		return fmt.Errorf("error walking through path: %v", err)
	}

	// sort.Slice(newFiless, func(i, j int) bool {
	// 	return newFiless[i].size < newFiless[j].size
	// })

	// for _, v := range newFiless {
	// 	queue.Enqueue(v.name)
	// }

	totalBytes := 0
	buf := *new(bytes.Buffer)
	zipWriter := zip.NewWriter(&buf)
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
		// buf.Reset()
		// if err = zipWriter.Flush(); err != nil {
		// 	return err
		// }

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
		fmt.Printf("%v consuming %v\n", Goid(), zipTask)
		if err := archiver.Archive(zipTask.zipFileList, zipTask.dest); err != nil {
			log.Fatal(err)
		}
	}
}
