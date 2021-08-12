package ttyimg

import (
	"image"
	"image/color"
	"math/bits"
	"sort"
)

func isqrt(x uint32) uint32 {
	r := x
	if r > 0xffff {
		r = 0xffff
	}
	l := uint32(0)

	for l+1 < r {
		c := (l + r) / 2
		v := c * c
		if v == x {
			return c
		}
		if v > x {
			r = c
		} else {
			l = c
		}
	}

	return l
}

func rgbaDist(r1, g1, b1, r2, g2, b2 uint32) uint32 {
	return isqrt((r1-r2)*(r1-r2) + (g1-g2)*(g1-g2) + (b1-b2)*(b1-b2))
}

func rgbaColDist(r1, g1, b1 uint32, c color.RGBA) uint32 {
	return rgbaDist(r1, g1, b1, uint32(c.R), uint32(c.G), uint32(c.B))
}

func colDist(c1, c2 color.RGBA) uint32 {
	return rgbaDist(uint32(c1.R), uint32(c1.G), uint32(c1.B), uint32(c2.R), uint32(c2.G), uint32(c2.B))
}

func generatePalette(cols []color.RGBA) (color.RGBA, color.RGBA, uint64) {
	type yuvInfo struct {
		Index int
		Y     uint8
		U     uint8
		V     uint8
	}

	yuvs := make([]yuvInfo, 0, len(cols))
	for i, c := range cols {
		y, u, v := color.RGBToYCbCr(c.R, c.G, c.B)
		yuvs = append(yuvs, yuvInfo{i, y, u, v})
	}

	sort.SliceStable(yuvs, func(i, j int) bool {
		if yuvs[i].Y < yuvs[j].Y {
			return true
		}
		if yuvs[i].Y > yuvs[j].Y {
			return false
		}

		if yuvs[i].U < yuvs[j].U {
			return true
		}
		if yuvs[i].U > yuvs[j].U {
			return false
		}

		return yuvs[i].V < yuvs[j].V
	})

	half := len(cols) / 2
	y1 := 0
	u1 := 0
	v1 := 0
	cnt1 := 0
	y2 := 0
	u2 := 0
	v2 := 0
	cnt2 := 0
	for i := 0; i < len(cols); i++ {
		if i < half {
			y1 += int(yuvs[i].Y)
			u1 += int(yuvs[i].U)
			v1 += int(yuvs[i].V)
			cnt1++
		} else {
			y2 += int(yuvs[i].Y)
			u2 += int(yuvs[i].U)
			v2 += int(yuvs[i].V)
			cnt2++
		}
	}

	if cnt1 != 0 {
		y1 /= cnt1
		u1 /= cnt1
		v1 /= cnt1
	}

	if cnt2 != 0 {
		y2 /= cnt2
		u2 /= cnt2
		v2 /= cnt2
	}

	r1, g1, b1 := color.YCbCrToRGB(uint8(y1), uint8(u1), uint8(v1))
	r2, g2, b2 := color.YCbCrToRGB(uint8(y2), uint8(u2), uint8(v2))

	mask := uint64(0)
	for i, c := range yuvs {
		if c.Index < 64 && i < half {
			mask |= uint64(1) << c.Index
		}
	}

	return color.RGBA{uint8(r1), uint8(g1), uint8(b1), 255}, color.RGBA{uint8(r2), uint8(g2), uint8(b2), 255}, mask
}

func encodeAsTextPixel8x8(img image.Image, offsetX, offsetY int) TextImagePixel {
	data := make([]color.RGBA, 0, 64)

	for y := offsetY; y < offsetY+8; y++ {
		for x := offsetX; x < offsetX+8; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			data = append(data, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
		}
	}

	fg, bg, sample := generatePalette(data)

	// Searching for the nearest member:
	var pixel *TextImagePixel = nil
	bestDist := 0
	for r, v := range drawFont {
		dist := bits.OnesCount64(sample ^ v)
		if pixel == nil || dist < bestDist {
			bestDist = dist
			pixel = &TextImagePixel{r, fg, bg}
		}

		dist = bits.OnesCount64(sample ^ (v ^ 0))
		if dist < bestDist {
			bestDist = dist
			pixel = &TextImagePixel{r, bg, fg}
		}
	}

	return *pixel
}
