package spritesheet

import (
	"colour"
	"compositor"
	"image"
	"sprite"
	"testing"
	"utils/imageutils"
)

func TestGetSpritesheets(t *testing.T) {

	def := Definition{
		Palette:    colour.Palette{Entries: []colour.PaletteEntry{{R: 0, G: 0, B: 0}, {R: 255, G: 255, B: 255}}},
		Scale:      1.0,
		NumSprites: 2,
	}

	sheets := GetSpritesheets(def)
	testSpritesheet(t, sheets, "32bpp")
	testSpritesheet(t, sheets, "8bpp")
	testSpritesheet(t, sheets, "mask")
}

func testSpritesheet(t *testing.T, sheets Spritesheets, bpp string) {
	sheet, ok := sheets.Data[bpp]

	if !ok {
		t.Fatalf("no " + bpp + "spritesheet present in result")
	}

	expectedRect := image.Rectangle{Max: image.Point{X: spriteSpacing * 2, Y: totalHeight}}
	spriteRect := getTestSpriteRectangle(0, 1.0)
	expectedImg := getTestSpriteImage(spriteRect)

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
	return getSpriteSizeForAngle(angle, scale)
}

func getTestSpriteImage(rect image.Rectangle) image.Image {
	spr := sprite.GetUniformSprite(rect)
	img := image.NewRGBA(rect)
	compositor.Composite32bpp(spr, img, image.Point{}, rect)
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
		{90, 32, 26},
		{135, 26, 26},
		{180, 24, 26},
		{225, 26, 26},
		{270, 32, 26},
		{315, 26, 26},
	}

	for _, testCase := range testCases {
		rect := getSpriteSizeForAngle(testCase.angle, 1.0)
		if rect.Max.X != testCase.expectedX || rect.Max.Y != testCase.expectedY {
			t.Errorf("output for angle %d was [%d,%d] (expected [%d,%d]", testCase.angle, rect.Max.X, rect.Max.Y, testCase.expectedX, testCase.expectedY)
		}
	}
}
