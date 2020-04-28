package spritesheet

import (
	"colour"
	"compositor"
	"image"
	"image/color"
	"sprite"
	"utils/imageutils"
	"voxelobject"
)

type Spritesheet struct {
	Image image.Image
}

type Spritesheets map[string]Spritesheet

const spriteSpacing = 40
const totalHeight = 40

func GetSpritesheets(object voxelobject.RawVoxelObject, pal colour.Palette, scale float64, numSprites int) Spritesheets {
	w := int(float64(spriteSpacing*numSprites) * scale)
	h := int(float64(totalHeight) * scale)
	bounds := image.Rectangle{Max: image.Point{X: w, Y: h}}

	img := imageutils.GetUniformImage(bounds, color.White)
	angleStep := 360 / numSprites

	for i := 0; i < numSprites; i++ {
		angle := 180 - (i * angleStep)
		sw, sh := getSpriteSizeForAngle(angle, scale)
		rect := image.Rectangle{Max: image.Point{X: sw, Y: sh}}
		var spr image.Image
		if object.Invalid() {
			spr = sprite.GetUniformSprite(rect)
		} else {
			spr = sprite.GetRaycastSprite(object, pal, rect, angle)
		}
		compositor.Composite(spr, img, image.Point{X: int(float64(i * spriteSpacing) * scale)}, rect)
	}

	sheets := make(Spritesheets)
	sheets["32bpp"] = Spritesheet{Image: img}

	return sheets
}

func getSpriteSizeForAngle(angle int, scale float64) (x, y int) {
	var fx, fy float64

	switch {
	case angle == 0 || angle == 180:
		fx, fy = 24, 26
	case angle == 90 || angle == 270:
		fx, fy = 32, 24
	default:
		fx, fy = 26, 26
	}

	return int(fx * scale), int(fy * scale)
}
