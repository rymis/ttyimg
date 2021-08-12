package ttyimg

import "image"

// Render image based on used fonts and TextImage representation
func RenderTextImage(ti *TextImage) *image.RGBA {
	res := image.NewRGBA(image.Rect(0, 0, ti.Width*8, ti.Height*16))

	for y := 0; y < ti.Height; y++ {
		for x := 0; x < ti.Width; x++ {
			px := ti.Get(x, y)
			mask, ok := drawFont[px.Char]
			one := uint64(1) << 63
			if ok {
				for yy := 0; yy < 8; yy++ {
					for xx := 0; xx < 8; xx++ {
						if (mask & one) != 0 {
							res.SetRGBA(x*8+xx, y*16+yy*2, px.Color)
							res.SetRGBA(x*8+xx, y*16+yy*2+1, px.Color)
						} else {
							res.SetRGBA(x*8+xx, y*16+yy*2, px.BackgroundColor)
							res.SetRGBA(x*8+xx, y*16+yy*2+1, px.BackgroundColor)
						}
						mask <<= 1
					}
				}
			}
		}
	}

	return res
}
