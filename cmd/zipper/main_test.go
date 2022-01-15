package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/mholt/archiver/v3"
	"github.com/stretchr/testify/assert"
	"gopkg.in/loremipsum.v1"
)

func TestE2E(t *testing.T) {
	// ch := make(chan struct{}, runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	noOfTmpFiles := 20

	// Create temp dir for creating test files
	pwd, err := os.Getwd()
	if err != nil {
		t.Error("Error getting working dir:", err)
	}
	testDir := filepath.Join(pwd, "testdir")
	tmpFilesPath := filepath.Join(testDir, "testfiles")
	if err = os.MkdirAll(tmpFilesPath, 0777); err != nil {
		t.Error("Error creating ", tmpFilesPath, " dir:", err)
	}

	// Create testfiles
	for i := 0; i < noOfTmpFiles; i++ {
		// tmpFileName := filepath.Join(tmpFilesPath, fmt.Sprint(i)+"-tmpfile.txt")
		// ch <- struct{}{}
		func(tmpFileName string) {
			testFile, err := os.Create(tmpFileName)
			if err != nil {
				t.Error("Error while creating file at ", tmpFileName, " :", err)
			}
			testFile.Write([]byte(loremipsum.New().Sentences(10000)))
			testFile.Close()
			// <-ch
		}(filepath.Join(tmpFilesPath, fmt.Sprint(i)+"-tmpfile.txt"))
	}

	// Create zip Test
	if err = crateZips(tmpFilesPath, 3000000); err != nil {
		t.Error("Error creating zip from path", tmpFilesPath, ": ", err)
	}

	// Remove tmp files, not zips
	defer func() {
		if err = os.RemoveAll(testDir); err != nil {
			t.Error("Error removing ", tmpFilesPath, ":", err)
		}
	}()

	// Extract files
	extractedZips := filepath.Join(pwd, "extractedzips")
	if err = os.MkdirAll(extractedZips, 0777); err != nil {
		t.Error("Error creating ", extractedZips, " dir:", err)
	}
	if err = filepath.Walk(pwd,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(info.Name()) == ".zip" {
				if err = Unarchive(filepath.Join(pwd, info.Name()), extractedZips); err != nil {
					t.Error("Error while Unarchiving", info.Name(), ":", err)
				}
				if err = os.Remove(filepath.Join(pwd, info.Name())); err != nil {
					t.Error("Error while removing ", filepath.Join(pwd, info.Name()), ":", err)
				}
			}
			return nil
		}); err != nil {
		t.Error("Error while walking through", pwd, ":", err)
	}

	// Count should be equal to no of files created
	count := 0
	if err = filepath.Walk(extractedZips,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(info.Name()) == ".txt" {
				count++
			}
			return nil
		}); err != nil {
		t.Error("Error while walking through", extractedZips, ":", err)
	}
	// Remove extracted files
	defer func() {
		if err = os.RemoveAll(extractedZips); err != nil {
			t.Error("Error removing ", extractedZips, ":", err)
		}
	}()
	assert.Equal(t, count, noOfTmpFiles, "Extracted files are ", fmt.Sprint(count), " and No of temp files are", noOfTmpFiles)
}
