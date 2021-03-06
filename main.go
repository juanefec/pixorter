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

	"github.com/juanefec/pixorter/sorter"

	"fmt"
	"image"
)

func main() {
	fileName := flag.String("filename", "image2.jpg", "use like: -filename=image.jpeg")
	flag.Parse()
	// You can register another format here
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	file := ReadImageFile(*fileName)
	sorted := sorter.Sort(file)
	WriteImageFile(sorted, *fileName)

}

func ReadImageFile(fileName string) *os.File {
	file, err := os.Open("./" + fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return file
}

func WriteImageFile(img image.Image, fname string) {
	name := "SORTED" + fname
	// open a new file
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}

	// and encoding it
	fmt := filepath.Ext(name)
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
