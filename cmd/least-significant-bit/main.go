package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/image/bmp"
)

const (
	byteSize = 8
)

func openBMP(path string) (image.Image, image.Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, image.Config{}, err
	}

	reader := bytes.NewReader(file)

	img, err := bmp.Decode(reader)
	if err != nil {
		return nil, image.Config{}, err
	}

	reader = bytes.NewReader(file)

	imgConfig, err := bmp.DecodeConfig(reader)
	if err != nil {
		return nil, image.Config{}, err
	}

	return img, imgConfig, nil
}

func getEncodedBytes(data []byte, message string, terminateSymbol rune) ([]byte, error) {
	msg := fmt.Sprintf("%s%s", message, string(terminateSymbol))

	if len(msg)*byteSize > len(data) {
		return nil, fmt.Errorf("message is to long, need %v bytes", len(msg)*byteSize)
	}

	encoded := make([]byte, len(data))
	copy(encoded, data)

	for i, char := range []byte(msg) {
		for j := 0; j < byteSize; j++ {
			idx := i*byteSize + j
			bit := char & (1 << j)

			encoded[idx] = (encoded[idx] &^ 1) | (bit >> j)
		}
	}

	return encoded, nil
}

func getDecodedMessage(data []byte, terminateSymbol rune) string {
	maxLength := len(data) / byteSize
	msgBytes := make([]byte, 0)

	for i := 0; i < maxLength; i++ {
		var char byte

		for j := 0; j < byteSize; j++ {
			bit := data[i*byteSize+j] & 1
			char = char | (bit << j)
		}

		if char == byte(terminateSymbol) {
			break
		}

		msgBytes = append(msgBytes, char)
	}

	return string(msgBytes)
}

func cloneImage(path string) (*image.NRGBA, error) {
	img, conf, err := openBMP(path)
	if err != nil {
		return nil, err
	}

	width := conf.Width
	height := conf.Height

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	newImg := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}

	return newImg, nil
}

func encodeMessageToImage(sourcePath, resultPath, message string) error {
	img, err := cloneImage(sourcePath)
	if err != nil {
		return err
	}

	img.Pix, err = getEncodedBytes(img.Pix, message, '#')
	if err != nil {
		return err
	}

	file, err := os.Create(resultPath)
	if err != nil {
		return err
	}

	return bmp.Encode(file, img)
}

func decodeMessageFromImage(path string) (string, error) {
	img, _, err := openBMP(path)
	if err != nil {
		return "", err
	}

	encoded := img.(*image.NRGBA)
	message := getDecodedMessage(encoded.Pix, '#')

	return message, nil
}

func main() {
	sourcePath := "images/samples/VENUS.BMP"
	resultPath := "images/encoded/result.bmp"
	message := "Hello, world!"

	err := encodeMessageToImage(sourcePath, resultPath, message)
	if err != nil {
		log.Fatal(err)
	}

	decodedMessage, err := decodeMessageFromImage(resultPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded message is: '%s'\n", decodedMessage)
}
