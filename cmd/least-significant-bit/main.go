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

func getEncodedBytes(data []byte, message string, terminateSymbol rune, period, offset int) ([]byte, error) {
	msg := fmt.Sprintf("%s%s", message, string(terminateSymbol))

	if len(msg)*byteSize*period > len(data) {
		return nil, fmt.Errorf("message is to long, need %v bytes", len(msg)*byteSize)
	}

	if period == 0 || offset > period {
		return nil, fmt.Errorf("period must be greater than 0 and greater than index")
	}

	encoded := make([]byte, len(data))
	copy(encoded, data)

	for i, char := range []byte(msg) {
		for j := 0; j < byteSize; j++ {
			k := i*byteSize + j
			idx := k*period + offset

			bit := char & (1 << j)

			encoded[idx] = (encoded[idx] &^ 1) | (bit >> j)
		}
	}

	return encoded, nil
}

func getDecodedMessage(data []byte, terminateSymbol rune, period, offset int) string {
	maxLength := len(data) / (byteSize * period)
	msgBytes := make([]byte, 0)

	for i := 0; i < maxLength; i++ {
		var char byte

		for j := 0; j < byteSize; j++ {
			k := i*byteSize + j
			idx := k*period + offset

			bit := data[idx] & 1
			char = char | (bit << j)
		}

		if char == byte(terminateSymbol) {
			break
		}

		msgBytes = append(msgBytes, char)
	}

	return string(msgBytes)
}

func cloneImage(path string) (*image.RGBA, error) {
	img, conf, err := openBMP(path)
	if err != nil {
		return nil, err
	}

	width := conf.Width
	height := conf.Height

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	newImg := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}

	return newImg, nil
}

func encodeMessageToImage(sourcePath, resultPath, message string, period, offset int) error {
	img, err := cloneImage(sourcePath)
	if err != nil {
		return err
	}

	img.Pix, err = getEncodedBytes(img.Pix, message, '#', period, offset)
	if err != nil {
		return err
	}

	file, err := os.Create(resultPath)
	if err != nil {
		return err
	}

	return bmp.Encode(file, img)
}

func decodeMessageFromImage(path string, period, offset int) (string, error) {
	img, _, err := openBMP(path)
	if err != nil {
		return "", err
	}

	encoded := img.(*image.RGBA)
	message := getDecodedMessage(encoded.Pix, '#', period, offset)

	return message, nil
}

func main() {
	sourcePath := "images/samples/VENUS.BMP"
	resultPath := "images/encoded/result.bmp"

	message := "Hello, world!"

	// 4 компоненты - R, G, B, A
	period := 4
	// для кодирования выбираем компоненту G
	offset := 1

	err := encodeMessageToImage(sourcePath, resultPath, message, period, offset)
	if err != nil {
		log.Fatal(err)
	}

	decodedMessage, err := decodeMessageFromImage(resultPath, period, offset)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded message is: '%s'\n", decodedMessage)
}
