package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mholt/archiver/v3"
	"gopkg.in/loremipsum.v1"
)

func TestZipFile(t *testing.T) {
	// Create tempdir
	tempDir, err := ioutil.TempDir("", "tmp")
	if err != nil {
		t.Errorf("Error while creating temp dir: %v", err)
	}

	// Give tmpFile
	tmpFileName := "zippertest.txt"
	tmpFilePath := filepath.Join(tempDir, tmpFileName)
	destZipName := filepath.Join(tempDir, tmpFileName) + ".zip"

	// Create file
	testFile, err := os.Create(tmpFilePath)
	if err != nil {
		t.Errorf("Error while creating file at \"%v\": %v\n", tmpFilePath, err)
	}
	testFile.Close()

	// Zip file
	err = archiver.Archive([]string{tmpFilePath}, destZipName)
	if err != nil {
		t.Errorf("Got error while zipping \"%v\" to \"%v\": %v\n", tmpFilePath, destZipName, err)
	}

	// Check file exist
	_, err = os.Stat(destZipName)
	if os.IsNotExist(err) {
		t.Errorf("Zip does not exist at \"%v\"\n", destZipName)
	}

	// Remove tmpFile and created zip file
	if err = os.RemoveAll(tmpFilePath); err != nil {
		t.Errorf("Error removing file: \"%v\"\n", tmpFilePath)
	}
	if err = os.RemoveAll(destZipName); err != nil {
		t.Errorf("Error removing file: \"%v\"\n", destZipName)
	}
}

func TestCopyFile(t *testing.T) {
	// Create tempdir
	tempDir1, err := ioutil.TempDir("", "tmp-")
	if err != nil {
		t.Errorf("Error while creating tempDir1: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir1); err != nil {
			t.Errorf("Error removing file: \"%v\"\n", tempDir1)
		}
	}()
	tmpFileName := "zippertest.txt"
	tmpFilePath := filepath.Join(tempDir1, tmpFileName)

	tempDir2, err := ioutil.TempDir("", "tmp-")
	if err != nil {
		t.Errorf("Error while creating tempDir2: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir2); err != nil {
			t.Errorf("Error removing file: \"%v\"\n", tempDir2)
		}
	}()

	// Create file
	testFile, err := os.Create(tmpFilePath)
	if err != nil {
		t.Errorf("Error while creating file at \"%v\": %v\n",
			tmpFilePath, err)
	}
	testFile.Write([]byte(loremipsum.New().Sentences(10000)))
	testFile.Close()

	// Copy tmpFile to tempDir2
	err = copy(tmpFilePath, filepath.Join(tempDir2, tmpFileName))
	if err != nil {
		t.Errorf("Cannot copy \"%v\" to \"%v\"\n",
			tmpFilePath, filepath.Join(tempDir2, tmpFileName))
	}

	// Check file exist
	_, err = os.Stat(filepath.Join(tempDir2, tmpFileName))
	if os.IsNotExist(err) {
		t.Errorf("file does not exist at \"%v\"\n",
			filepath.Join(tempDir2, tmpFileName))
	}
}
