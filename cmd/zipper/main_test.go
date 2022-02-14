package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	. "github.com/mholt/archiver/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"gopkg.in/loremipsum.v1"
)

func TestE2E(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// noOfTmpFiles := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(50)
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
		func(tmpFileName string) {
			testFile, err := os.Create(tmpFileName)
			if err != nil {
				t.Error("Error while creating file at ", tmpFileName, " :", err)
			}
			// testFile.Write([]byte(loremipsum.New().Sentences(4000)))
			testFile.Write([]byte(loremipsum.New().Sentences(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(4000))))
			testFile.Close()
		}(filepath.Join(tmpFilesPath, fmt.Sprint(i)+"-tmpfile.txt"))
	}

	// Create zip Test
	if err = crateZips(tmpFilesPath, 3000000); err != nil {
		t.Error("Error creating zip from path", tmpFilesPath, ": ", err)
	}

	// Check goroutine leak test
	goleak.VerifyNone(t)

	// Remove tmp files, not zips
	t.Cleanup(func() {
		if err = os.RemoveAll(testDir); err != nil {
			t.Error("Error removing ", tmpFilesPath, ":", err)
		}
	})

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
				t.Cleanup(func() {
					if err = os.Remove(filepath.Join(pwd, info.Name())); err != nil {
						t.Error("Error while removing ", filepath.Join(pwd, info.Name()), ":", err)
					}
				})
			}
			return nil
		}); err != nil {
		t.Error("Error while walking through", pwd, ":", err)
	}

	// Count should be equal to no of files created
	count := 0
	noOfExtractedFiles, err := ioutil.ReadDir(extractedZips + osPathSuffix)
	if err != nil {
		panic(err)
	}

	// Remove extracted files
	assert.Equal(t, len(noOfExtractedFiles), noOfTmpFiles, "Extracted files are ", fmt.Sprint(count), " and No of temp files are", noOfTmpFiles)
	t.Cleanup(func() {
		if err = os.RemoveAll(extractedZips); err != nil {
			t.Error("Error removing ", extractedZips, ":", err)
		}
	})
	time.Sleep(time.Second * 2)
}

// goleak can also be run at the end of every test package by creating a
// TestMain function for your package:
func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
