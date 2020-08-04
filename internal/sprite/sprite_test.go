package sprite

import (
	"github.com/mattkimber/gorender/internal/utils/imageutils"
	"image"
	"image/color"
	"testing"
)

func TestApplySprite(t *testing.T) {
	rect := image.Rectangle{Max: image.Point{X: 2, Y: 2}}
	img := imageutils.GetUniformImage(rect, color.White)
	ApplyUniformSprite(img, rect, image.Point{})

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
