package sorter

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Pixel struct example

type img struct {
	h, w int
	m    [][]color.RGBA
}
type Px color.RGBA

func (m *img) At(x, y int) color.Color { return m.m[x][y] }
func (m *img) ColorModel() color.Model { return color.RGBAModel }
func (m *img) Bounds() image.Rectangle { return image.Rect(0, 0, m.h, m.w) }

// Sort wraps the sorting process
func Sort(imgfile io.ReadCloser) image.Image {
	rand.Seed(time.Now().Unix())
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

	var pixels [][]color.RGBA
	for x := 0; x < width; x++ {
		var row []color.RGBA
		for y := 0; y < height; y++ {
			row = append(row, rgbaToRGBA(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}
	fmt.Println(height, width)
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
	workers := 10
	var wg sync.WaitGroup
	wg.Add(workers + 1)
	c := make(chan struct{ i, j int })
	pch := make(chan pxm)
	for i := 0; i < workers; i++ {
		go func() {
			for t := range c {
				fillPixel(m, n, t.i, t.j, pch)
			}
			wg.Done()
		}()
	}
	
	chosen := make([]pxm, 0)
	go func() {
		for p := range pch {
			chosen = append(chosen, p)
		}
		wg.Done()
	}()

	for i, row := range m.m {
		for j := range row {
			c <- struct{ i, j int }{i, j}
		}
	}

	close(c)
	close(pch)
	wg.Wait()
	//fmt.Println(chosen)
	for i := range chosen {
		cant := random(-25, 25)
		for w := 0; w < cant; w++ {
			safeRefillPixel(n, chosen[i].c, chosen[i].i+w, chosen[i].j)
		}
	}
}

func fillPixel(m img, n *img, i, j int, pch chan pxm) {
	//dis(m, n, i, j)
	//datdaat(m, n, i, j)
	juat(m, n, i, j, pch)

}

func safeRefillPixel(n *img, c color.RGBA, i, j int) {
	canI := i < len(n.m) && i > 0
	canJ := j < len(n.m[0]) && j > 0
	if canI && canJ {
		refillPixel(n, c, i, j)
	}
}

func refillPixel(n *img, c color.RGBA, i, j int) {
	n.m[i][j].R, n.m[i][j].G, n.m[i][j].B, n.m[i][j].A = 170, 170, 170, 170
}

type pxm struct {
	i, j int
	c    color.RGBA
}
