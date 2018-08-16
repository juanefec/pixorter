package sorter

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"sync"
)

// Pixel struct example

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

	originalImg := img{w, h, pixels}
	newImg := make([][]color.RGBA, len(pixels))
	for i := range newImg {
		newImg[i] = make([]color.RGBA, len(pixels[i]))
	}
	outImg := &img{w, h, newImg}
	workersFillImg(originalImg, outImg)
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
func workersFillImg(m img, n *img) {
	workers := 32
	var wg sync.WaitGroup
	wg.Add(workers)
	c := make(chan struct{ i, j int })
	for i := 0; i < workers; i++ {
		go func() {
			for t := range c {
				fillPixel(m, n, t.i, t.j)
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

func fillPixel(m img, n *img, i, j int) {
	//dis(m, n, i, j)
	dat(m, n, i, j)

}
func dis(m img, n *img, i, j int) {
	if i > 15 && j > 13 {
		ni, nj := i-15, j-13
		n.m[i][j].R, n.m[i][j].G, n.m[i][j].B, n.m[i][j].A = m.m[ni][nj].R, m.m[ni][nj].G, m.m[ni][nj].B, m.m[ni][nj].A

	}
}
func dat(m img, n *img, i, j int) {
	if rgbAvg(m.m[i][j]) {
		nR := m.m[i][j].R - 20
		nG := m.m[i][j].G + 90
		nB := m.m[i][j].B - 30
		n.m[i][j].R, n.m[i][j].G, n.m[i][j].B, n.m[i][j].A = nR, nG, nB, m.m[i][j].A
	} else {
		n.m[i][j] = m.m[i][j]
	}
}

func rgbAvg(c color.RGBA) bool {
	avg := int((c.R + c.G + c.B + c.A) / 4)
	dr, dg, db := abs(uint8(avg)-c.R), abs(uint8(avg)-c.G), abs(uint8(avg)-c.B)
	return dr < 120 && dg < 130 && db < 90
}

func abs(n uint8) uint8 {
	return uint8(math.Abs(float64(n)))
}
