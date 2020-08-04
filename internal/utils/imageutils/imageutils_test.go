package imageutils

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"testing"
)

func TestGetUniformImage(t *testing.T) {
	rect := image.Rectangle{Max: image.Point{X: 2, Y: 2}}
	img := GetUniformImage(rect, color.Black)

	if img.Bounds() != rect {
		t.Errorf("Image bounds %v not equal to expected %v", img.Bounds(), rect)
	}

	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()

			if r != 0 || g != 0 || b != 0 {
				t.Errorf("Non-black pixel %v at %d,%d", img.At(x, y), x, y)
			}
		}
	}
}

func TestIsColourEqual(t *testing.T) {
	rect := image.Rectangle{Max: image.Point{X: 2, Y: 2}}
	img := GetUniformImage(rect, color.Black)

	if !IsColourEqual(img, 1, 1, 0, 0, 0) {
		t.Errorf("Expected colour to be equal but isn't")
	}

	if IsColourEqual(img, 1, 1, 255, 255, 255) {
		t.Errorf("Expected colour to not be equal but is")
	}
}

func TestIsImageEqualToSubImage(t *testing.T) {
	rect1 := image.Rectangle{Max: image.Point{X: 3, Y: 3}}
	rect2 := image.Rectangle{Max: image.Point{X: 1, Y: 1}}
	img1 := GetUniformImage(rect1, color.Black)
	img2 := GetUniformImage(rect2, color.White)

	draw.Draw(img1, rect2, img2, image.Point{}, draw.Src)

	if !IsImageEqualToSubImage(img2, img1, rect2) {
		t.Errorf("Expected equality at %v but was not equal", rect2)
	}

	if IsImageEqualToSubImage(img2, img1, rect2.Add(image.Point{X: 1})) {
		t.Errorf("Expected equality at %v but was not equal", rect2.Add(image.Point{X: 1}))
	}
}

func TestClearToColourIndex(t *testing.T) {
	const expected = 5
	rect := image.Rectangle{Max: image.Point{X: 3, Y: 3}}
	img := image.NewPaletted(rect, palette.Plan9)
	ClearToColourIndex(img, expected)

	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			if img.ColorIndexAt(x, y) != expected {
				t.Errorf("colour at %d %d is %d - expected %d", x, y, img.ColorIndexAt(x, y), expected)
			}
		}
	}
}
