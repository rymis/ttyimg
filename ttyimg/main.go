package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"github.com/rymis/ttyimg"
	"image"
	imgcolor "image/color"
	"image/jpeg"
	_ "image/png"
	"os"
)

func check(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printImage16Color(img *ttyimg.TextImage) {
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			p := img.Get(x, y)
			fmt.Printf("%c", p.Char)
		}
		fmt.Println()
	}
}

func colIndex256(c imgcolor.RGBA, bg imgcolor.RGBA) *color.Style256 {
	// Actual used values for colors are:
	// 0, 95, 135, 175, 215, 255
	// My interval borders are:
	// [0..47] [48..115] [116..155] [156..195] [196..235] [236..255]
	convert := func(c uint8) uint32 {
		if c <= 47 {
			return 0
		}
		if c <= 115 {
			return 1
		}
		if c <= 155 {
			return 2
		}
		if c <= 195 {
			return 3
		}
		if c <= 235 {
			return 4
		}
		return 5
	}

	r := convert(c.R)
	g := convert(c.G)
	b := convert(c.B)

	r2 := convert(bg.R)
	g2 := convert(bg.G)
	b2 := convert(bg.B)

	return color.S256(uint8(16+r*36+g*6+b), uint8(16+r2*36+g2*6+b2))
}

func printImage256Color(img *ttyimg.TextImage) {
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			p := img.Get(x, y)
			s := colIndex256(p.Color, p.BackgroundColor)
			s.Printf("%c", p.Char)
		}
		fmt.Println()
	}
}

func printImageNoColor(img *ttyimg.TextImage) {
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			p := img.Get(x, y)
			fmt.Printf("%c", p.Char)
		}
		fmt.Println()
	}
}

func main() {
	input := flag.String("input", "", "input image file to process as tty image")
	output := flag.String("output", "", "output image to render image into (if present only)")
	width := flag.Int("width", 80, "text image width")
	flag.Parse()

	if *input == "" {
		fmt.Printf("Error: input argument is mandatory")
		os.Exit(1)
	}

	if *width <= 0 {
		*width = 80
	}

	imgFile, err := os.Open(*input)
	check(err)
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	check(err)

	txtimg := ttyimg.NewTextImageFromImage(img, *width)
	if color.IsSupport256Color() {
		printImage256Color(txtimg)
	} else if color.IsSupportColor() {
		printImage16Color(txtimg)
	} else {
		printImageNoColor(txtimg)
	}

	if output != nil && *output != "" {
		imgout := ttyimg.RenderTextImage(txtimg)
		f, err := os.Create(*output)
		if err != nil {
			fmt.Printf("Error: can't open file %s\n", *output)
			os.Exit(1)
		}
		defer f.Close()

		jpeg.Encode(f, imgout, nil)
	}
}
