package sorter

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"sync"
)

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}

type img struct {
	h, w int
	m    [][]color.RGBA
}

func (m *img) At(x, y int) color.Color { return m.m[x][y] }
func (m *img) ColorModel() color.Model { return color.RGBAModel }
func (m *img) Bounds() image.Rectangle { return image.Rect(0, 0, m.h, m.w) }

func Sort(imgfile io.ReadCloser) image.Image {
	pixels, h, w, err := GetPixels(imgfile)
	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}

	outImg := &img{w, h, pixels}

	workersFillImg(outImg)
	return outImg
}

// GetPixels gets the bi-dimensional pixel array
func GetPixels(imgfile io.ReadCloser) ([][]color.RGBA, int, int, error) {
	defer imgfile.Close()
	img, _, err := image.Decode(imgfile)
	if err != nil {
		return nil, 0, 0, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	fmt.Println(height, width)
	var pixels [][]color.RGBA
	for x := 0; x < width; x++ {
		var row []color.RGBA
		for y := 0; y < height; y++ {
			row = append(row, rgbaToRGBA(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, height, width, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToRGBA(r uint32, g uint32, b uint32, a uint32) color.RGBA {
	return color.RGBA{uint8(r / 257), uint8(g / 257), uint8(b / 257), uint8(a / 257)}
}

// 4 workers per CPU
// real	0m17.304s
// user	0m40.615s
// sys	0m2.517s
func workersFillImg(m *img) {
	var workers int = 3000
	var wg sync.WaitGroup
	wg.Add(workers)

	c := make(chan struct{ i, j int })
	for i := 0; i < workers; i++ {
		go func() {
			for t := range c {
				fillPixel(m, t.i, t.j)
			}
			wg.Done()
		}()
	}

	for i, row := range m.m {
		for j := range row {
			c <- struct{ i, j int }{i, j}
		}
	}
	close(c)
	wg.Wait()
}

func fillPixel(m *img, i, j int) {
	m.m[i][j].R, m.m[i][j].G, m.m[i][j].B = m.m[i][j].G, m.m[i][j].B, m.m[i][j].R
}
