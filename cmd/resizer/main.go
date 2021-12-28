package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/disintegration/imaging"
)

func init() {
	h := sha1.New()
	h.Write([]byte(os.Getenv("UNLOCK_ZIPPER")))

	if fmt.Sprintf("%x", h.Sum(nil)) != "f43dc14eb92b5dafba73481449f8ed3a602c0b6e" {
		os.Exit(2)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Open a test imapath stringge.
	// imageFile, err := os.Open(path)
	dirPath := os.Args[1]

	// Reading images from dir
	var imageFiles []string
	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			// Append files in string arr
			imageFiles = append(
				imageFiles,
				filepath.Join(dirPath, info.Name()))
			return nil
		})
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range imageFiles {
		decodeImage(v)
	}
}

func decodeImage(imageFiles string) {
	// Decoding images
	imageFile, err := os.Open(imageFiles)
	if err != nil {
		log.Fatalf("Cannot open file: %v\n", err)
	}

	image, _, err := image.DecodeConfig(imageFile)
	if err != nil {
		log.Fatalf("Cannot decode config: %v\n", err)
	}
	imageFile.Close()

	log.Printf("Processing image: %v\n", imageFiles)
	switch {
	case image.Width > image.Height:
		if image.Width == 1920 && image.Height == 1080 {
			log.Printf("%v>%v: Lookslike landscape image\n",
				image.Width, image.Height)
		} else {
			log.Println("Resizing to landscape (1920x1080)")
			resizeImage(imageFiles, 1920, 1080)
		}
	case image.Width < image.Height:
		if image.Width == 1080 && image.Height == 1920 {
			log.Printf("%v>%v: Lookslike portrait image\n",
				image.Width, image.Height)
		} else {
			log.Println("Resizing to portrait (1080x1920)")
			resizeImage(imageFiles, 1080, 1920)
		}
	default:
		log.Println("Not compatible")
	}
}

func resizeImage(path string, width, height int) error {
	file, err := imaging.Open(path)
	if err != nil {
		log.Fatalf("Failed to open image \"%v\": %v\n", path, err)
		// return errors.New("Failed to open image \"%v\": %v\n", path, err)
	}
	// if err = imaging.Save(imaging.Resize(file, width, height, imaging.Lanczos), filepath.Join(filepath.Dir(path), strings.TrimRight(filepath.Base(path), filepath.Ext(path))+"-replace"+filepath.Ext(path))); err != nil {
	if err = os.Remove(path); err != nil {
		return errors.New("Failed to remove: " + path)
	}
	if err = imaging.Save(imaging.Resize(file, width, height, imaging.Lanczos), path); err != nil {
		log.Fatalf("Failed to save image \"%v\": %v\n", path, err)
	}
	return nil
}
