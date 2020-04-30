package compositor

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

func Composite(src image.Image, dst image.Image, loc image.Point, size image.Rectangle) error {
	writableDst, ok := dst.(draw.Image)
	if !ok {
		return fmt.Errorf("could not convert destination image to writable image")
	}

	w, h := size.Bounds().Max.X, size.Bounds().Max.Y
	xFactor, yFactor := float64(src.Bounds().Max.X)/float64(w), float64(src.Bounds().Max.Y)/float64(h)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			sx, sy := int(float64(x)*xFactor), int(float64(y)*yFactor)
			sw, sh := int(float64(x+1)*xFactor), int(float64(y+1)*yFactor)
			ct := 0
			r, g, b, a := 0, 0, 0, 0

			for i := sx; i < sw; i++ {
				for j := sy; j < sh; j++ {
					cr, cg, cb, ca := src.At(i, j).RGBA()
					r += int(cr)
					g += int(cg)
					b += int(cb)
					a += int(ca)
					ct++
				}
			}

			if ct == 0 {
				ct = 1
			}

			c := color.RGBA64{R: uint16(r / ct), G: uint16(g / ct), B: uint16(b / ct), A: uint16(a / ct)}
			writableDst.Set(x+loc.X, y+loc.Y, c)
		}
	}

	return nil
}
