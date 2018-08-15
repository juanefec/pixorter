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

var fileName = flag.String("filename", "image.jpg", "use like: -filename=image.jpeg")

func main() {
	flag.Parse()
	// You can register another format here
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	file, err := os.Open("./" + *fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sorted := sorter.Sort(file)
	SaveImageFile(sorted, *fileName)

}

func SaveImageFile(img image.Image, fname string) {
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
