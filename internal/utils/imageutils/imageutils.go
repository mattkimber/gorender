package imageutils

import (
	"image"
	"image/color"
	"image/draw"
)

func GetUniformImage(bounds image.Rectangle, colour color.Color) *image.RGBA {
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, &image.Uniform{C: colour}, image.Point{}, draw.Src)
	return img
}

func ClearToColourIndex(img *image.Paletted, index byte) {
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			img.SetColorIndex(x, y, index)
		}
	}
}

func IsColourEqual(img image.Image, x int, y int, r uint32, g uint32, b uint32) bool {
	ir, ig, ib, _ := img.At(x, y).RGBA()
	if ir != r || ig != g || ib != b {
		return false
	}
	return true
}

func IsImageEqualToSubImage(img image.Image, sub image.Image, bounds image.Rectangle) bool {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			sx, sy := x-bounds.Min.X, y-bounds.Min.Y

			r, g, b, a := img.At(x, y).RGBA()
			rs, gs, bs, as := sub.At(sx, sy).RGBA()

			if r != rs || g != gs || b != bs || a != as {
				return false
			}
		}
	}
	return true
}
