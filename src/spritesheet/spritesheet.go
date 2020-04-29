package spritesheet

import (
	"colour"
	"compositor"
	"image"
	"image/color"
	"image/png"
	"io"
	"sprite"
	"utils/fileutils"
	"utils/imageutils"
	"voxelobject"
)

type Spritesheet struct {
	Image image.Image
}

type Definition struct {
	Object     voxelobject.RawVoxelObject
	Palette    colour.Palette
	Scale      float64
	NumSprites int
}

type Spritesheets map[string]Spritesheet

const spriteSpacing = 40
const totalHeight = 40

func GetSpritesheets(def Definition) Spritesheets {
	sheets := make(Spritesheets)

	w := int(float64(spriteSpacing*def.NumSprites) * def.Scale)
	h := int(float64(totalHeight) * def.Scale)
	bounds := image.Rectangle{Max: image.Point{X: w, Y: h}}

	sheets["32bpp"] = Spritesheet{Image: getSpritesheetImage(def, bounds)}

	return sheets
}

func getSpritesheetImage(def Definition, bounds image.Rectangle) (img image.Image) {
	img = imageutils.GetUniformImage(bounds, color.White)
	angleStep := 360 / float64(def.NumSprites)

	for i := 0; i < def.NumSprites; i++ {
		angle := 180 - int(float64(i)*angleStep)
		spr := getSprite(def, angle)
		compositor.Composite(spr, img, image.Point{X: int(float64(i*spriteSpacing) * def.Scale)}, spr.Bounds())
	}

	return
}

func getSprite(def Definition, angle int) (spr image.Image) {
	rect := getSpriteSizeForAngle(angle, def.Scale)

	if def.Object.Invalid() {
		spr = sprite.GetUniformSprite(rect)
	} else {
		spr = sprite.GetRaycastSprite(def.Object, def.Palette, rect, angle)
	}

	return
}

func (s Spritesheet) OutputToWriter(w io.Writer) (err error) {
	err = png.Encode(w, s.Image)
	return
}

func (sheets Spritesheets) SaveAll(baseFilename string) (err error) {
	for i, sheet := range sheets {
		filename := baseFilename + "_" + i + ".png"
		if err = fileutils.WriteToFile(filename, &sheet); err != nil {
			return
		}
	}

	return
}

func getSpriteSizeForAngle(angle int, scale float64) image.Rectangle {
	var fx, fy float64

	switch {
	case angle == 0 || angle == 180:
		fx, fy = 24, 26
	case angle == 90 || angle == 270:
		fx, fy = 32, 24
	default:
		fx, fy = 26, 26
	}

	return image.Rectangle{Max: image.Point{X: int(fx * scale), Y: int(fy * scale)}}
}
