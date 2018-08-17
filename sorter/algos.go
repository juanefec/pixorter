package sorter

import (
	"image/color"
	"math"
	"math/rand"
)

func dis(m img, n *img, i, j int) {
	if i > 15 && j > 13 {
		ni, nj := i-15, j-13
		n.m[i][j].R, n.m[i][j].G, n.m[i][j].B, n.m[i][j].A = m.m[ni][nj].R, m.m[ni][nj].G, m.m[ni][nj].B, m.m[ni][nj].A

	}
}
func datdaat(m img, n *img, i, j int) {
	if dat(m, n, i, j) && daat(m, n, i, j) {
		n.m[i][j] = m.m[i][j]
	}
}

func dat(m img, n *img, i, j int) bool {
	if rgbAvg(m.m[i][j]) {
		nR := m.m[i][j].R
		nG := m.m[i][j].G + 60
		nB := m.m[i][j].B - 30
		n.m[i][j].R, n.m[i][j].G, n.m[i][j].B, n.m[i][j].A = nR, nG, nB, m.m[i][j].A
		return false
	} else {
		//n.m[i][j] = m.m[i][j]
		return true
	}
}

func daat(m img, n *img, i, j int) bool {
	ni, nj := i+random(-2, 2), j+random(-2, 2)
	//fmt.Println(ni, nj)
	canI := ni < len(m.m) && ni > 0
	canJ := nj < len(m.m[i]) && nj > 0
	//fmt.Println(canI, canJ)
	if canI && canJ {
		n.m[i][j] = m.m[ni][nj]
		return false
	} else {
		//n.m[i][j] = m.m[i][j]
		return true
	}
}

func juat(m img, n *img, i, j int, pch chan pxm) {
	datdaat(m, n, i, j)
	c := m.m[i][j]
	if int((c.R+c.G+c.B)/3) > 75 {
		pch <- pxm{i, j, c}
	}
}

// _____________________________common_____________________________

func rgbAvg(c color.RGBA) bool {
	avg := int((c.R + c.G + c.B + c.A) / 4)
	dr, dg, db := abs(uint8(avg)-c.R), abs(uint8(avg)-c.G), abs(uint8(avg)-c.B)
	return dr < 120 && dg < 130 && db < 90
}

func abs(n uint8) uint8 {
	return uint8(math.Abs(float64(n)))
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
