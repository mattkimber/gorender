package compositor

import (
	"colour"
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

type DestinationImageSampler func(rect image.Rectangle, loc image.Point)

func Composite32bpp(src image.Image, dst image.Image, loc image.Point, size image.Rectangle) error {
	writableDst, ok := dst.(draw.Image)
	if !ok {
		return fmt.Errorf("could not convert destination image to writable image")
	}

	sampler := func(rect image.Rectangle, pt image.Point) {
		c := resample32bpp(src, rect)
		writableDst.Set(pt.X, pt.Y, c)
	}

	return composite(src.Bounds(), sampler, loc, size)
}

func Composite8bpp(src *image.Paletted, dst *image.Paletted, loc image.Point, size image.Rectangle, pal colour.Palette) error {
	sampler := func(rect image.Rectangle, pt image.Point) {
		c := resample8bpp(src, rect, pal)
		dst.SetColorIndex(pt.X, pt.Y, c)
	}

	return composite(src.Bounds(), sampler, loc, size)

}

func composite(srcBounds image.Rectangle, dst DestinationImageSampler, loc image.Point, size image.Rectangle) error {
	w, h := size.Bounds().Max.X, size.Bounds().Max.Y
	xFactor, yFactor := float64(srcBounds.Max.X)/float64(w), float64(srcBounds.Max.Y)/float64(h)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			rect := image.Rectangle{
				Min: image.Point{X: int(float64(x) * xFactor), Y: int(float64(y) * yFactor)},
				Max: image.Point{X: int(float64(x+1) * xFactor), Y: int(float64(y+1) * yFactor)},
			}
			dst(rect, image.Point{X: x + loc.X, Y: y + loc.Y})
		}
	}

	return nil
}

func resample32bpp(src image.Image, bounds image.Rectangle) color.RGBA64 {
	ct := 0
	r, g, b, a := 0, 0, 0, 0

	for i := bounds.Min.X; i < bounds.Max.X; i++ {
		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
			cr, cg, cb, ca := src.At(i, j).RGBA()
			r += int(cr)
			g += int(cg)
			b += int(cb)
			a += int(ca)
			ct++
		}
	}

	return color.RGBA64{R: uint16(r / ct), G: uint16(g / ct), B: uint16(b / ct), A: uint16(a / ct)}
}

func resample8bpp(src *image.Paletted, bounds image.Rectangle, pal colour.Palette) byte {
	values := map[byte]int{}

	for i := bounds.Min.X; i < bounds.Max.X; i++ {
		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
			idx := src.ColorIndexAt(i, j)

			if pal.IsCompanyColour(idx) {
				return idx
			}
			values[idx]++
		}
	}

	max := 0
	modalIndex := byte(0)

	for k, v := range values {
		if v > max {
			max = v
			modalIndex = k
		}
	}

	return modalIndex
}
