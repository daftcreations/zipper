package main

import (
	"image"
	_ "image/jpeg"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLandscapImage(t *testing.T) {
	imagePath := `../../misc/landscap-image.png`

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

	decodeImage(imagePath)

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
	imagePath := `../../misc/portait-image.png`

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

	decodeImage(imagePath)

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
