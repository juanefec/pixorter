package main

import (
	"errors"
	"flag"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"fmt"
	"image"
	"io"
)

func main() {
	// You can register another format here
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	file, err := os.Open("./image.jpg")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	pixels, err := getPixels(file)

	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}

	fmt.Println(pixels)
}

// Get the bi-dimensional pixel array
func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}

var (
	output  = flag.String("out", "mandelbrot.png", "name of the output image file")
	height  = flag.Int("h", 2048, "height of the output image in pixels")
	width   = flag.Int("w", 2048, "width of the output image in pixels")
	mode    = flag.String("mode", "seq", "mode: seq, px, row, workers")
	workers = flag.Int("workers", 1, "number of workers to use")
)

func SaveImageFile() {
	flag.Parse()

	// open a new file
	f, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}

	// create the image
	var img image.Image

	// and encoding it
	fmt := filepath.Ext(*output)
	switch fmt {
	case ".png":
		err = png.Encode(f, img)
	case ".jpg", ".jpeg":
		err = jpeg.Encode(f, img, nil)
	case ".gif":
		err = gif.Encode(f, img, nil)
	default:
		err = errors.New("unkwnown format " + fmt)
	}
	// unless you can't
	if err != nil {
		log.Fatal(err)
	}
}
