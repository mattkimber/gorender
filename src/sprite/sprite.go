package sprite

import (
	"colour"
	"image"
	"image/color"
)

func ApplyUniformSprite(img *image.RGBA, bounds image.Rectangle, loc image.Point) {
	minX, minY := bounds.Min.X+loc.X, bounds.Min.Y+loc.Y
	maxX, maxY := bounds.Max.X+loc.X, bounds.Max.Y+loc.Y

	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			img.Set(x, y, color.Black)
		}
	}
}

func Apply32bppSprite(img *image.RGBA, bounds image.Rectangle, loc image.Point, info ShaderOutput, getProperty func(*ShaderInfo) colour.RGB) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			c := getProperty(&info[x][y])
			img.Set(x+loc.X, y+loc.Y, c.GetRGBA(info[x][y].Alpha))
		}
	}
}

func ApplyIndexedSprite(img *image.Paletted, bounds image.Rectangle, loc image.Point, info ShaderOutput, getProperty func(*ShaderInfo) byte) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			c := getProperty(&info[x][y])
			img.SetColorIndex(x+loc.X, y+loc.Y, c)
		}
	}
}
