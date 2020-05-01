package sprite

import (
	"colour"
	"geometry"
	"image"
	"image/color"
	"raycaster"
	"utils/imageutils"
)

type shadeFunc32bpp func(int, int) color.RGBA64
type shadeFuncIndexed func(int, int) byte

func GetUniformSprite(bounds image.Rectangle) image.Image {
	return imageutils.GetUniformImage(bounds, color.Black)
}

func get32bppImage(bounds image.Rectangle, shader shadeFunc32bpp, info raycaster.RenderOutput) image.Image {
	img := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if info[x][y].Collision {
				img.Set(x, y, shader(x, y))
			}
		}
	}

	return img
}

func getIndexedImage(pal colour.Palette, bounds image.Rectangle, shader shadeFuncIndexed, info raycaster.RenderOutput) *image.Paletted {
	img := image.NewPaletted(bounds, pal.GetGoPalette())

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if info[x][y].Collision {
				img.SetColorIndex(x, y, shader(x, y))
			}
		}
	}

	return img
}

func Get32bppSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) image.Image {
	shader := func(x, y int) color.RGBA64 {
		lightingOffset := (info[x][y].LightAmount * 0.6) +
			((-(float64(info[x][y].Depth - 120) / 40)) * 0.1) - 0.2
		r, g, b := pal.GetLitRGB(info[x][y].Index, lightingOffset)
		return color.RGBA64{R: r, G: g, B: b, A: 65535}
	}

	return get32bppImage(bounds, shader, info)
}

func GetIndexedSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) *image.Paletted {
	shader := func(x, y int) byte {
		return info[x][y].Index
	}

	return getIndexedImage(pal, bounds, shader, info)
}

func GetMaskSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) *image.Paletted {
	shader := func(x, y int) byte {
		return pal.GetMaskColour(info[x][y].Index)
	}

	return getIndexedImage(pal, bounds, shader, info)
}

func GetNormalSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) image.Image {
	shader := func(x, y int) color.RGBA64 {
		normal := info[x][y].Normal.MultiplyByConstant(32766).Add(geometry.Vector3{X: 32766, Y: 32766, Z: 32766})
		r, g, b := normal.X, normal.Y, normal.Z
		return color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: 65535}
	}

	return get32bppImage(bounds, shader, info)
}

func GetAverageNormalSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) image.Image {
	shader := func(x, y int) color.RGBA64 {
		normal := info[x][y].AveragedNormal.MultiplyByConstant(32766).Add(geometry.Vector3{X: 32766, Y: 32766, Z: 32766})
		r, g, b := normal.X, normal.Y, normal.Z
		return color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: 65535}
	}

	return get32bppImage(bounds, shader, info)
}

func GetDepthSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) image.Image {
	shader := func(x, y int) color.RGBA64 {
		v := info[x][y].Depth * 400
		return color.RGBA64{R: uint16(v), G: uint16(v), B: uint16(v), A: 65535}
	}

	return get32bppImage(bounds, shader, info)
}

func GetLightingSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) image.Image {
	shader := func(x, y int) color.RGBA64 {
		v := 32767 + (info[x][y].LightAmount * 32767)
		return color.RGBA64{R: uint16(v), G: uint16(v), B: uint16(v), A: 65535}
	}

	return get32bppImage(bounds, shader, info)
}
