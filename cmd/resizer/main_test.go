package main

import (
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLandscapImage(t *testing.T) {
	imagePath := `../../misc/to-landscap-image.png`

	// Original image
	imageFile, err := os.Open(imagePath)
	if err != nil {
		t.Errorf("Cannot open file: %v\n", err)
	}
	decodedImage, _, err := image.DecodeConfig(imageFile)
	if err != nil {
		t.Errorf("Cannot decode config: %v\n", err)
	}
	imageFile.Close()

	if err := decodeImage(imagePath); err != nil {
		t.Errorf("Cannot decode config: %v\n", err)
	}

	// Resize image
	resizedImage, err := os.Open(imagePath)
	if err != nil {
		t.Errorf("Cannot open file: %v\n", err)
	}
	resizedDecodeImage, _, err := image.DecodeConfig(resizedImage)
	if err != nil {
		t.Errorf("Cannot decode config: %v\n", err)
	}
	resizedImage.Close()

	assert.NotEqual(t,
		decodedImage.Width, resizedDecodeImage.Width,
		"New images' Width shouldn't be equal")
	assert.NotEqual(t,
		decodedImage.Height, resizedDecodeImage.Height,
		"New images' Height shouldn't be equal")
}

func TestPotraitImage(t *testing.T) {
	imagePath := `../../misc/to-portait-image.png`

	// Original image
	imageFile, err := os.Open(imagePath)
	if err != nil {
		t.Errorf("Cannot open file: %v\n", err)
	}
	decodedImage, _, err := image.DecodeConfig(imageFile)
	if err != nil {
		t.Errorf("Cannot decode config: %v\n", err)
	}
	imageFile.Close()

	if err := decodeImage(imagePath); err != nil {
		log.Fatalln(err)
	}

	// Resize image
	resizedImage, err := os.Open(imagePath)
	if err != nil {
		t.Errorf("Cannot open file: %v\n", err)
	}
	resizedDecodeImage, _, err := image.DecodeConfig(resizedImage)
	if err != nil {
		t.Errorf("Cannot decode config: %v\n", err)
	}
	resizedImage.Close()

	assert.NotEqual(t,
		decodedImage.Width, resizedDecodeImage.Width,
		"New images' Width shouldn't be equal")
	assert.NotEqual(t,
		decodedImage.Height, resizedDecodeImage.Height,
		"New images' Height shouldn't be equal")
}

func TestLandscapImage2(t *testing.T) {
	imagePath := `../../misc/landscap-image.png`

	if err := decodeImage(imagePath); err != nil {
		t.Errorf("Cannot decode config: %v\n", err)
	}
}

func TestNonImage(t *testing.T) {
	imagePath := `../../misc/non-image.txt`

	if !assert.NotNil(t, decodeImage(imagePath), "okayy") {
		t.Error("Should return error like \"image: unknown format\"")
	}
}
