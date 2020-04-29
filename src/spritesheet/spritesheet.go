package spritesheet

import (
	"colour"
	"compositor"
	"image"
	"image/color"
	"image/png"
	"io"
	"raycaster"
	"sprite"
	"utils/fileutils"
	"utils/imageutils"
	"voxelobject"
)

type Spritesheet struct {
	Image image.Image
}

type Definition struct {
	Object     voxelobject.ProcessedVoxelObject
	Palette    colour.Palette
	Scale      float64
	NumSprites int
	Debug      bool
}

type Spritesheets map[string]Spritesheet

type SpriteInfo struct {
	RenderOutput raycaster.RenderOutput
	Bounds       image.Rectangle
}

const spriteSpacing = 40
const totalHeight = 40

func GetSpritesheets(def Definition) Spritesheets {
	sheets := make(Spritesheets)

	w := int(float64(spriteSpacing*def.NumSprites) * def.Scale)
	h := int(float64(totalHeight) * def.Scale)
	bounds := image.Rectangle{Max: image.Point{X: w, Y: h}}
	spriteInfos := make([]SpriteInfo, def.NumSprites)

	angleStep := 360 / float64(def.NumSprites)
	for i := 0; i < def.NumSprites; i++ {
		angle := ((180 - int(float64(i)*angleStep)) + 360) % 360
		rect := getSpriteSizeForAngle(angle, def.Scale)

		spriteInfos[i].RenderOutput = raycaster.GetRaycastOutput(def.Object, angle, rect.Max.X, rect.Max.Y, def.Debug)
		spriteInfos[i].Bounds = rect
	}

	sheets["32bpp"] = Spritesheet{Image: get32bppSpritesheetImage(def, bounds, spriteInfos, "32bpp")}
	sheets["8bpp"] = Spritesheet{Image: get8bppSpritesheetImage(def, bounds, spriteInfos, "8bpp")}
	sheets["mask"] = Spritesheet{Image: get8bppSpritesheetImage(def, bounds, spriteInfos, "mask")}

	if def.Debug {
		sheets["lighting"] = Spritesheet{Image: get32bppSpritesheetImage(def, bounds, spriteInfos, "lighting")}
		sheets["depth"] = Spritesheet{Image: get32bppSpritesheetImage(def, bounds, spriteInfos, "depth")}
		sheets["normals"] = Spritesheet{Image: get32bppSpritesheetImage(def, bounds, spriteInfos, "normal")}
		sheets["avg_normals"] = Spritesheet{Image: get32bppSpritesheetImage(def, bounds, spriteInfos, "avg")}
	}

	return sheets
}

func get8bppSpritesheetImage(def Definition, bounds image.Rectangle, spriteInfos []SpriteInfo, depth string) image.Image {
	palette := def.Palette.GetGoPalette()
	img := image.NewPaletted(bounds, palette)
	imageutils.ClearToColourIndex(img, byte(len(palette)-1))

	for i := 0; i < def.NumSprites; i++ {
		spr := getSprite(def, spriteInfos[i], depth)
		compositor.Composite(spr, img, image.Point{X: int(float64(i*spriteSpacing) * def.Scale)}, spr.Bounds())
	}

	return img
}

func get32bppSpritesheetImage(def Definition, bounds image.Rectangle, spriteInfos []SpriteInfo, depth string) (img image.Image) {
	img = imageutils.GetUniformImage(bounds, color.White)

	for i := 0; i < def.NumSprites; i++ {
		spr := getSprite(def, spriteInfos[i], depth)
		compositor.Composite(spr, img, image.Point{X: int(float64(i*spriteSpacing) * def.Scale)}, spr.Bounds())
	}

	return
}

func getSprite(def Definition, spriteInfo SpriteInfo, depth string) (spr image.Image) {
	if def.Object.Invalid() {
		spr = sprite.GetUniformSprite(spriteInfo.Bounds)
	} else if depth == "8bpp" {
		spr = sprite.GetIndexedSprite(def.Palette, spriteInfo.Bounds, spriteInfo.RenderOutput)
	} else if depth == "mask" {
		spr = sprite.GetMaskSprite(def.Palette, spriteInfo.Bounds, spriteInfo.RenderOutput)
	} else if depth == "lighting" {
		spr = sprite.GetLightingSprite(def.Palette, spriteInfo.Bounds, spriteInfo.RenderOutput)
	}  else if depth == "depth" {
		spr = sprite.GetDepthSprite(def.Palette, spriteInfo.Bounds, spriteInfo.RenderOutput)
	} else if depth == "normal" {
		spr = sprite.GetNormalSprite(def.Palette, spriteInfo.Bounds, spriteInfo.RenderOutput)
	}  else if depth == "avg" {
		spr = sprite.GetAverageNormalSprite(def.Palette, spriteInfo.Bounds, spriteInfo.RenderOutput)
	} else {
		spr = sprite.Get32bppSprite(def.Palette, spriteInfo.Bounds, spriteInfo.RenderOutput)
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
