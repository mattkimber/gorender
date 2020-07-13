package spritesheet

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"manifest"
	"raycaster"
	"sampler"
	"sprite"
	"sync"
	"utils/fileutils"
	"utils/imageutils"
	timingutils "utils/timingutils"
)

type Spritesheet struct {
	Image image.Image
}

type Spritesheets struct {
	sync.RWMutex
	Data map[string]Spritesheet
}

type SpriteInfo struct {
	ShaderOutput sprite.ShaderOutput
	SpriteBounds image.Rectangle
}

const spriteSpacing = 8

func GetSpritesheets(def manifest.Definition) (sheets Spritesheets) {
	sheets.Data = make(map[string]Spritesheet)

	w, h := 0, 0
	for i, spr := range def.Manifest.Sprites {
		def.Manifest.Sprites[i].X = w
		w += int(float64(spr.Width+spriteSpacing) * def.Scale)

		if int(float64(spr.Height)*def.Scale) > h {
			h = int(float64(spr.Height) * def.Scale)
		}
	}

	bounds := image.Rectangle{Max: image.Point{X: w, Y: h}}
	spriteInfos := make([]SpriteInfo, len(def.Manifest.Sprites))

	timingutils.Time("Raycasting/sampling", def.Time, func() {
		raycast(def, spriteInfos)
	})

	timingutils.Time("Spritesheets", def.Time, func() {
		getRegularSheets(&sheets, def, bounds, spriteInfos)
	})
	if def.Debug {
		timingutils.Time("Debug output", def.Time, func() {
			getDebugSheets(&sheets, def, bounds, spriteInfos)
		})
	}

	return
}

func getDebugSheets(sheets *Spritesheets, def manifest.Definition, bounds image.Rectangle, spriteInfos []SpriteInfo) {
	debugOutputs := []string{"lighting", "depth", "normals", "occlusion", "shadow", "avg_normals"}
	var wg sync.WaitGroup
	wg.Add(len(debugOutputs) + 1)

	for _, s := range debugOutputs {
		thisS := s
		go func() {
			defer wg.Done()
			sheets.Store(thisS, Spritesheet{Image: get32bppSpritesheetImage(def, bounds, spriteInfos, thisS)})
		}()
	}

	go func() {
		defer wg.Done()
		smp := sampler.Get(def.Manifest.Sampler)(1, 1, def.Manifest.Accuracy, def.Manifest.Overlap)
		sheets.Store("sampler", Spritesheet{Image: smp.GetImage()})
	}()

	wg.Wait()
}

func getRegularSheets(sheets *Spritesheets, def manifest.Definition, bounds image.Rectangle, spriteInfos []SpriteInfo) {
	var wg sync.WaitGroup
	wg.Add(1)
	if !def.Only8bpp {
		wg.Add(2)
	}

	go func() {
		defer wg.Done()
		sheets.Store("8bpp", Spritesheet{Image: get8bppSpritesheetImage(def, bounds, spriteInfos, "8bpp")})
	}()

	if !def.Only8bpp {
		go func() {
			defer wg.Done()
			sheets.Store("32bpp", Spritesheet{Image: get32bppSpritesheetImage(def, bounds, spriteInfos, "32bpp")})
		}()
		go func() {
			defer wg.Done()
			sheets.Store("mask", Spritesheet{Image: get8bppSpritesheetImage(def, bounds, spriteInfos, "mask")})
		}()
	}

	wg.Wait()
}

func raycast(def manifest.Definition, spriteInfos []SpriteInfo) {
	for i, spr := range def.Manifest.Sprites {
		rect := getSpriteSizeForAngle(spr, def.Scale)

		smpFunc := sampler.Get(def.Manifest.Sampler)
		smp := smpFunc(rect.Max.X, rect.Max.Y, def.Manifest.Accuracy, def.Manifest.Overlap)

		spriteInfos[i].SpriteBounds = rect
		renderOutput := raycaster.GetRaycastOutput(def.Object, def.Manifest, spr, smp)
		spriteInfos[i].ShaderOutput = sprite.GetShaderOutput(renderOutput, def, rect.Max.X, rect.Max.Y)
	}
}

func get8bppSpritesheetImage(def manifest.Definition, bounds image.Rectangle, spriteInfos []SpriteInfo, depth string) image.Image {
	palette := def.Palette.GetGoPalette()
	img := image.NewPaletted(bounds, palette)
	imageutils.ClearToColourIndex(img, byte(len(palette)-1))

	for i := 0; i < len(def.Manifest.Sprites); i++ {
		loc := image.Point{X: def.Manifest.Sprites[i].X}
		applySprite8bpp(img, spriteInfos[i], loc, depth)
	}

	return img
}

func get32bppSpritesheetImage(def manifest.Definition, bounds image.Rectangle, spriteInfos []SpriteInfo, depth string) image.Image {
	img := imageutils.GetUniformImage(bounds, color.White)

	for i := 0; i < len(def.Manifest.Sprites); i++ {
		loc := image.Point{X: def.Manifest.Sprites[i].X}
		applySprite32bpp(img, def, spriteInfos[i], loc, depth)
	}

	return img
}

func applySprite8bpp(img *image.Paletted, spriteInfo SpriteInfo, loc image.Point, depth string) {
	if depth == "8bpp" {
		sprite.ApplyIndexedSprite(img, spriteInfo.SpriteBounds, loc, spriteInfo.ShaderOutput, sprite.GetIndex)
	} else if depth == "mask" {
		sprite.ApplyIndexedSprite(img, spriteInfo.SpriteBounds, loc, spriteInfo.ShaderOutput, sprite.GetMaskIndex)
	}

	return
}

func applySprite32bpp(img *image.RGBA, def manifest.Definition, spriteInfo SpriteInfo, loc image.Point, depth string) {
	if def.Object.Invalid() {
		sprite.ApplyUniformSprite(img, spriteInfo.SpriteBounds, loc)
	} else if depth == "lighting" {
		sprite.Apply32bppSprite(img, spriteInfo.SpriteBounds, loc, spriteInfo.ShaderOutput, sprite.GetLighting)
	} else if depth == "depth" {
		sprite.Apply32bppSprite(img, spriteInfo.SpriteBounds, loc, spriteInfo.ShaderOutput, sprite.GetDepth)
	} else if depth == "occlusion" {
		sprite.Apply32bppSprite(img, spriteInfo.SpriteBounds, loc, spriteInfo.ShaderOutput, sprite.GetOcclusion)
	} else if depth == "shadow" {
		sprite.Apply32bppSprite(img, spriteInfo.SpriteBounds, loc, spriteInfo.ShaderOutput, sprite.GetShadowing)
	} else if depth == "normals" {
		sprite.Apply32bppSprite(img, spriteInfo.SpriteBounds, loc, spriteInfo.ShaderOutput, sprite.GetNormal)
	} else if depth == "avg_normals" {
		sprite.Apply32bppSprite(img, spriteInfo.SpriteBounds, loc, spriteInfo.ShaderOutput, sprite.GetAveragedNormal)
	} else {
		sprite.Apply32bppSprite(img, spriteInfo.SpriteBounds, loc, spriteInfo.ShaderOutput, sprite.GetColour)
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

func (sheets *Spritesheets) SaveAll(baseFilename string) (err error) {
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

func getSpriteSizeForAngle(sprite manifest.Sprite, scale float64) image.Rectangle {
	fx, fy := float64(sprite.Width), float64(sprite.Height)
	return image.Rectangle{Max: image.Point{X: int(fx * scale), Y: int(fy * scale)}}
}
