package spritesheet

import (
	"compositor"
	"image"
	"sprite"
	"testing"
	"utils/imageutils"
	"voxelobject"
)

func TestGetSpritesheets(t *testing.T) {
	expectedRect := image.Rectangle{Max: image.Point{X: spriteSpacing * 2, Y: totalHeight}}
	spriteRect := getTestSpriteRectangle(0, 1.0)
	expectedImg := getTestSpriteImage(spriteRect)

	sheets := GetSpritesheets(voxelobject.RawVoxelObject{}, 1.0, 2)
	sheet, ok := sheets["32bpp"]

	if !ok {
		t.Fatalf("no 32bpp spritesheet present in result")
	}

	if sheet.Image.Bounds() != expectedRect {
		t.Errorf("spritesheet size %v did not match expected size %v", sheet.Image.Bounds(), expectedRect)
	}

	if !imageutils.IsImageEqualToSubImage(sheet.Image, expectedImg, spriteRect) {
		t.Errorf("sprite at %v not equal to composited output", spriteRect)
	}

	if !imageutils.IsColourEqual(sheet.Image, spriteSpacing-1, 0, 65535, 65535, 65535) {
		t.Errorf("blank area of spritesheet not set to white")
	}
}

func getTestSpriteRectangle(angle int, scale float64) image.Rectangle {
	x, y := getSpriteSizeForAngle(angle, scale)
	rect := image.Rectangle{Max: image.Point{X: x, Y: y}}
	return rect
}

func getTestSpriteImage(rect image.Rectangle) image.Image {
	spr := sprite.GetSprite(rect, 0)
	img := image.NewRGBA(rect)
	compositor.Composite(spr, img, image.Point{}, rect)
	return img
}

func TestGetSpriteSizeForAngle(t *testing.T) {
	testCases := []struct {
		angle     int
		expectedX int
		expectedY int
	}{
		{0, 24, 26},
		{45, 26, 26},
		{90, 32, 24},
		{135, 26, 26},
		{180, 24, 26},
		{225, 26, 26},
		{270, 32, 24},
		{315, 26, 26},
	}

	for _, testCase := range testCases {
		x, y := getSpriteSizeForAngle(testCase.angle, 1.0)
		if x != testCase.expectedX || y != testCase.expectedY {
			t.Errorf("output for angle %d was [%d,%d] (expected [%d,%d]", testCase.angle, x, y, testCase.expectedX, testCase.expectedY)
		}
	}
}
