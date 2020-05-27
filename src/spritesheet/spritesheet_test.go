package spritesheet

import (
	"colour"
	"compositor"
	"geometry"
	"image"
	"manifest"
	"sprite"
	"testing"
	"utils/imageutils"
)

func TestGetSpritesheets(t *testing.T) {

	def := Definition{
		Palette: colour.Palette{Entries: []colour.PaletteEntry{{R: 0, G: 0, B: 0}, {R: 255, G: 255, B: 255}}},
		Scale:   1.0,
		Manifest: manifest.Manifest{
			LightingAngle:        45,
			LightingElevation:    60,
			Size:                 geometry.Vector3{},
			RenderElevationAngle: 0,
			Sprites: []manifest.Sprite{
				{Angle: 0, Width: 32, Height: 32, X: 0},
				{Angle: 45, Width: 32, Height: 32, X: 40},
			},
		},
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

	expectedRect := image.Rectangle{Max: image.Point{X: 80, Y: 32}}
	spriteRect := getTestSpriteRectangle(manifest.Sprite{Width: 32, Height: 32}, 1.0)
	expectedImg := getTestSpriteImage(spriteRect)

	if sheet.Image.Bounds() != expectedRect {
		t.Errorf("spritesheet size %v did not match expected size %v", sheet.Image.Bounds(), expectedRect)
	}

	if !imageutils.IsImageEqualToSubImage(sheet.Image, expectedImg, spriteRect) {
		t.Errorf("sprite at %v not equal to composited output", spriteRect)
	}

	if !imageutils.IsColourEqual(sheet.Image, 79, 0, 65535, 65535, 65535) {
		t.Errorf("blank area of spritesheet not set to white")
	}
}

func getTestSpriteRectangle(spr manifest.Sprite, scale float64) image.Rectangle {
	return getSpriteSizeForAngle(spr, scale)
}

func getTestSpriteImage(rect image.Rectangle) image.Image {
	spr := sprite.GetUniformSprite(rect)
	img := image.NewRGBA(rect)
	compositor.Composite32bpp(spr, img, image.Point{}, rect, manifest.Manifest{})
	return img
}
