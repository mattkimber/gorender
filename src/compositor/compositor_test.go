package compositor

import (
	"colour"
	"image"
	"image/color"
	"image/color/palette"
	"manifest"
	"testing"
	"utils/imageutils"
)

func TestComposite32bpp(t *testing.T) {
	r1 := image.Rectangle{Max: image.Point{X: 1, Y: 1}}
	r2 := image.Rectangle{Max: image.Point{X: 2, Y: 2}}
	img1 := imageutils.GetUniformImage(r1, color.Black)
	img2 := imageutils.GetUniformImage(r2, color.White)

	err := Composite32bpp(img1, img2, image.Point{X: 1}, r1, manifest.Manifest{})

	if err != nil {
		t.Errorf("could not convert image to writable format: %s", err)
	}

	testColorAt(img2, 0, 0, 65535, 65535, 65535, t)
	testColorAt(img2, 1, 0, 0, 0, 0, t)
	testColorAt(img2, 1, 1, 65535, 65535, 65535, t)
}

func TestComposite8bpp(t *testing.T) {
	r1 := image.Rectangle{Max: image.Point{X: 1, Y: 1}}
	r2 := image.Rectangle{Max: image.Point{X: 2, Y: 2}}
	img1 := imageutils.GetUniformPalettedImage(r1, palette.Plan9, 0)
	img2 := imageutils.GetUniformPalettedImage(r2, palette.Plan9, 255)

	err := Composite8bpp(img1, img2, image.Point{X: 1}, r1, colour.Palette{})

	if err != nil {
		t.Errorf("could not convert image to writable format: %s", err)
	}

	testColorAt(img2, 0, 0, 65535, 65535, 65535, t)
	testColorAt(img2, 1, 0, 0, 0, 0, t)
	testColorAt(img2, 1, 1, 65535, 65535, 65535, t)
}

func testColorAt(img image.Image, x int, y int, r uint32, g uint32, b uint32, t *testing.T) {
	if !imageutils.IsColourEqual(img, x, y, r, g, b) {
		t.Errorf("Pixel at [%d,%d] is not equal to [%d,%d,%d])", x, y, r, g, b)
	}
}
