package sprite

import (
	"image"
	"testing"
	"utils/imageutils"
)

func TestGetSprite(t *testing.T) {
	rect := image.Rectangle{Max: image.Point{X: 2, Y: 2}}
	img := GetSprite(rect, 0)

	if img.Bounds() != rect {
		t.Errorf("Image bounds %v not equal to expected %v", img.Bounds(), rect)
	}

	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			if !imageutils.IsColourEqual(img, x, y, 0, 0, 0) {
				t.Errorf("Non-black pixel %v at 1,1", img.At(1, 1))
			}
		}
	}
}
