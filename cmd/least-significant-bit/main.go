package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"

	"golang.org/x/image/bmp"
)

const (
	testMessage = "Hello, world!"
	byteSize    = 8
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
			encoded[idx] = encoded[idx] | (bit >> j)
		}
	}

	return encoded, nil
}

func getDecodedMessage(data []byte, terminateSymbol rune) string {
	fmt.Printf("%16b", terminateSymbol)

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

		fmt.Printf("%08b: %+v\n", char, char)

		if i == 10 {
			break
		}

		msgBytes = append(msgBytes, char)
	}

	return string(msgBytes)
}

func encodeMessage(message, terminateSymbol string, img image.Image) (image.Image, error) {
	return nil, nil
}

func decodeMessage(message string, img image.Image) (image.Image, error) {
	return nil, nil
}

func main() {
	img2, conf, err := openBMP("images/sample.bmp")
	fmt.Printf("%+v\n", conf)

	// fmt.Printf("%+v", imgConfig)

	// pixels := make([]byte, 100*100)

	// img := image.NewGray(image.Rect(0, 0, 100, 100))
	// img.Pix = pixels

	// img = &image.Gray{Pix: pixels, Stride: 100, Rect: image.Rect(0, 0, 100, 100)}

	// img = image.NewGray(image.Rect(0, 0, 100, 100))
	// copy(img.Pix, pixels)

	width := conf.Width
	height := conf.Height

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	// cyan := color.RGBA{100, 200, 200, 0xff}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, img2.At(x, y))
		}
	}

	a := make([]byte, 2000)

	// encoded, err := getEncodedBytes(img.Pix, "Test text!!!", '#')
	encoded, err := getEncodedBytes(a, "Test text!!!", '#')
	if err != nil {
		log.Fatal(err)
	}

	// img.Pix = encoded

	// f, _ := os.Create("images/encoded.bmp")
	// bmp.Encode(f, img)

	// fmt.Println(getDecodedMessage(img.Pix, '#'))
	fmt.Println(getDecodedMessage(encoded, '#'))

}
