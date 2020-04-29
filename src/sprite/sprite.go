package sprite

import (
	"colour"
	"image"
	"image/color"
	"raycaster"
	"utils/imageutils"
)

func GetUniformSprite(bounds image.Rectangle) image.Image {
	return imageutils.GetUniformImage(bounds, color.Black)
}

func Get32bppSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) image.Image {
	img := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if info[x][y].Collision {
				r, g, b := pal.GetRGB(info[x][y].Index)
				img.Set(x, y, color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: 65535})
			}
		}
	}

	return img
}

func GetIndexedSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) image.Image {
	img := image.NewPaletted(bounds, pal.GetGoPalette())

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if info[x][y].Collision {
				img.SetColorIndex(x, y, info[x][y].Index)
			}
		}
	}

	return img
}

func GetMaskSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) image.Image {
	img := image.NewPaletted(bounds, pal.GetGoPalette())

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if info[x][y].Collision {
				index := pal.GetMaskColour(info[x][y].Index)
				img.SetColorIndex(x, y, index)
			}
		}
	}

	return img
}
