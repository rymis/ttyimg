package ttyimg

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"
)

// Pixel information in textual image
type TextImagePixel struct {
	Char            rune
	Color           color.RGBA
	BackgroundColor color.RGBA
}

type TextImage struct {
	Width  int
	Height int
	Data   []TextImagePixel
}

func NewTextImage(width, height int) *TextImage {
	res := &TextImage{}
	res.Width = width
	res.Height = height
	res.Data = make([]TextImagePixel, 0, width*height)

	for i := 0; i < width*height; i++ {
		res.Data = append(res.Data, TextImagePixel{' ', color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}})
	}

	return res
}

func NewTextImageFromImage(img image.Image, width int) *TextImage {
	imgrect := img.Bounds()
	height := width * (imgrect.Max.Y - imgrect.Min.Y) / 2 / (imgrect.Max.X - imgrect.Min.X)
	rect := image.Rect(0, 0, width*8, height*8)
	small := image.NewRGBA(rect)

	draw.BiLinear.Scale(small, rect, img, img.Bounds(), draw.Over, nil)

	res := NewTextImage(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			res.SetPixel(x, y, encodeAsTextPixel8x8(small, x*8, y*8))
		}
	}

	return res
}

func (img *TextImage) Set(x, y int, ch rune, fg color.Color, bg color.Color) {
	if y < img.Height && x < img.Width {
		img.Data[y*img.Width+x] = TextImagePixel{ch, rgbaColor(fg), rgbaColor(bg)}
	}
}

func (img *TextImage) SetPixel(x, y int, p TextImagePixel) {
	if y < img.Height && x < img.Width {
		img.Data[y*img.Width+x] = p
	}
}

func (img *TextImage) Get(x, y int) TextImagePixel {
	if y < img.Height && x < img.Width {
		return img.Data[y*img.Width+x]
	}

	return TextImagePixel{' ', color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}}
}

func rgbaColor(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}
