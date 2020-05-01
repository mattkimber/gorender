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
	"sync"
	"utils/fileutils"
	"utils/imageutils"
	timeutils "utils/timingutils"
	"voxelobject"
)

type Spritesheet struct {
	Image image.Image
}

type Spritesheets struct {
	sync.RWMutex
	Data map[string]Spritesheet
}

type Definition struct {
	Object     voxelobject.ProcessedVoxelObject
	Palette    colour.Palette
	Scale      float64
	NumSprites int
	Debug      bool
	Time       bool
}

type SpriteInfo struct {
	RenderOutput raycaster.RenderOutput
	RenderBounds image.Rectangle
	SpriteBounds image.Rectangle
}

const spriteSpacing = 40
const totalHeight = 40
const antiAliasFactor = 2

func GetSpritesheets(def Definition) Spritesheets {
	sheets := Spritesheets{}
	sheets.Data = make(map[string]Spritesheet)

	w := int(float64(spriteSpacing*def.NumSprites) * def.Scale)
	h := int(float64(totalHeight) * def.Scale)
	bounds := image.Rectangle{Max: image.Point{X: w, Y: h}}
	spriteInfos := make([]SpriteInfo, def.NumSprites)

	timeutils.Time("Raycasting", def.Time, func() {
		raycast(def, spriteInfos)
	})

	timeutils.Time("Spritesheets", def.Time, func() {
		getRegularSheets(sheets, def, bounds, spriteInfos)
	})
	if def.Debug {
		timeutils.Time("Debug output", def.Time, func() {
			getDebugSheets(sheets, def, bounds, spriteInfos)
		})
	}

	return sheets
}

func getDebugSheets(sheets Spritesheets, def Definition, bounds image.Rectangle, spriteInfos []SpriteInfo) {
	debugOutputs := []string{"lighting", "depth", "normals", "occlusion", "shadow", "avg_normals"}
	var wg sync.WaitGroup
	wg.Add(len(debugOutputs))

	for _, s := range debugOutputs {
		thisS := s
		go func() {
			sheets.Store(thisS, Spritesheet{Image: get32bppSpritesheetImage(def, bounds, spriteInfos, thisS)})
			wg.Done()
		}()
	}

	wg.Wait()
}

func getRegularSheets(sheets Spritesheets, def Definition, bounds image.Rectangle, spriteInfos []SpriteInfo) {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		sheets.Store("32bpp", Spritesheet{Image: get32bppSpritesheetImage(def, bounds, spriteInfos, "32bpp")})
		wg.Done()
	}()
	go func() {
		sheets.Store("8bpp", Spritesheet{Image: get8bppSpritesheetImage(def, bounds, spriteInfos, "8bpp")})
		wg.Done()
	}()
	go func() {
		sheets.Store("mask", Spritesheet{Image: get8bppSpritesheetImage(def, bounds, spriteInfos, "mask")})
		wg.Done()
	}()

	wg.Wait()
}

func raycast(def Definition, spriteInfos []SpriteInfo) {
	angleStep := 360 / float64(def.NumSprites)
	for i := 0; i < def.NumSprites; i++ {
		angle := ((int(float64(i) * angleStep)) + 360) % 360
		rect := getSpriteSizeForAngle(angle, def.Scale)

		rw, rh := rect.Max.X*antiAliasFactor, rect.Max.Y*antiAliasFactor
		spriteInfos[i].SpriteBounds = rect
		spriteInfos[i].RenderBounds = image.Rectangle{Max: image.Point{X: rw, Y: rh}}
		spriteInfos[i].RenderOutput = raycaster.GetRaycastOutput(def.Object, angle, rw, rh)
	}
}

func get8bppSpritesheetImage(def Definition, bounds image.Rectangle, spriteInfos []SpriteInfo, depth string) image.Image {
	palette := def.Palette.GetGoPalette()
	img := image.NewPaletted(bounds, palette)
	imageutils.ClearToColourIndex(img, byte(len(palette)-1))

	for i := 0; i < def.NumSprites; i++ {
		spr := getSprite8bpp(def, spriteInfos[i], depth)
		compositor.Composite8bpp(spr, img, image.Point{X: int(float64(i*spriteSpacing) * def.Scale)}, spriteInfos[i].SpriteBounds, def.Palette)
	}

	return img
}

func get32bppSpritesheetImage(def Definition, bounds image.Rectangle, spriteInfos []SpriteInfo, depth string) (img image.Image) {
	img = imageutils.GetUniformImage(bounds, color.White)

	for i := 0; i < def.NumSprites; i++ {
		spr := getSprite32bpp(def, spriteInfos[i], depth)
		compositor.Composite32bpp(spr, img, image.Point{X: int(float64(i*spriteSpacing) * def.Scale)}, spriteInfos[i].SpriteBounds)
	}

	return
}

func getSprite8bpp(def Definition, spriteInfo SpriteInfo, depth string) (spr *image.Paletted) {
	if depth == "8bpp" {
		spr = sprite.GetIndexedSprite(def.Palette, spriteInfo.RenderBounds, spriteInfo.RenderOutput)
	} else if depth == "mask" {
		spr = sprite.GetMaskSprite(def.Palette, spriteInfo.RenderBounds, spriteInfo.RenderOutput)
	}

	return
}

func getSprite32bpp(def Definition, spriteInfo SpriteInfo, depth string) (spr image.Image) {
	if def.Object.Invalid() {
		spr = sprite.GetUniformSprite(spriteInfo.RenderBounds)
	} else if depth == "lighting" {
		spr = sprite.GetLightingSprite(def.Palette, spriteInfo.RenderBounds, spriteInfo.RenderOutput)
	} else if depth == "depth" {
		spr = sprite.GetDepthSprite(def.Palette, spriteInfo.RenderBounds, spriteInfo.RenderOutput)
	} else if depth == "occlusion" {
		spr = sprite.GetOcclusionSprite(def.Palette, spriteInfo.RenderBounds, spriteInfo.RenderOutput)
	} else if depth == "shadow" {
		spr = sprite.GetShadowSprite(def.Palette, spriteInfo.RenderBounds, spriteInfo.RenderOutput)
	} else if depth == "normal" {
		spr = sprite.GetNormalSprite(def.Palette, spriteInfo.RenderBounds, spriteInfo.RenderOutput)
	} else if depth == "avg" {
		spr = sprite.GetAverageNormalSprite(def.Palette, spriteInfo.RenderBounds, spriteInfo.RenderOutput)
	} else {
		spr = sprite.Get32bppSprite(def.Palette, spriteInfo.RenderBounds, spriteInfo.RenderOutput)
	}

	return
}

func (s Spritesheet) OutputToWriter(w io.Writer) (err error) {
	err = png.Encode(w, s.Image)
	return
}

func (sheets *Spritesheets) Store(key string, s Spritesheet) {
	sheets.Lock()
	sheets.Data[key] = s
	sheets.Unlock()
}

func (sheets Spritesheets) SaveAll(baseFilename string) (err error) {
	var wg sync.WaitGroup
	wg.Add(len(sheets.Data))

	for i, sheet := range sheets.Data {
		filename := baseFilename + "_" + i + ".png"
		thisSheet := sheet
		go func() { fileutils.WriteToFile(filename, thisSheet); wg.Done() }()
	}

	wg.Wait()
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
