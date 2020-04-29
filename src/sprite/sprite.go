package sprite

import (
	"colour"
	"geometry"
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
				v, cb, cr := color.RGBToYCbCr(byte(r>>8),byte(g>>8),byte(b>>8))

				v = byte((float64(v) * 0.4) +
					((1.0 + info[x][y].LightAmount) * (127.0 * 0.4)) +
					((float64(info[x][y].Depth / (340/127))) * 0.2))
				r2, g2, b2 := color.YCbCrToRGB(v, cb, cr)

				img.Set(x, y, color.RGBA{R: r2, G: g2, B: b2, A: 255})
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

func GetNormalSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) image.Image {
	img := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if info[x][y].Collision {
				normal := info[x][y].Normal.MultiplyByConstant(32766).Add(geometry.Vector3{X: 32766, Y: 32766, Z: 32766})
				r, g, b := normal.X, normal.Y, normal.Z
				img.Set(x, y, color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: 65535})
			}
		}
	}

	return img
}

func GetAverageNormalSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) image.Image {
	img := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if info[x][y].Collision {
				normal := info[x][y].AveragedNormal.MultiplyByConstant(32767).Add(geometry.Vector3{X: 32767, Y: 32767, Z: 32767})
				r, g, b := normal.X, normal.Y, normal.Z
				img.Set(x, y, color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: 65535})
			}
		}
	}

	return img
}

func GetDepthSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) image.Image {
	img := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if info[x][y].Collision {
				v := info[x][y].Depth * 100
				img.Set(x, y, color.RGBA64{R: uint16(v), G: uint16(v), B: uint16(v), A: 65535})
			}
		}
	}

	return img
}

func GetLightingSprite(pal colour.Palette, bounds image.Rectangle, info raycaster.RenderOutput) image.Image {
	img := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if info[x][y].Collision {
				v := 32767 + (info[x][y].LightAmount * 32767)
				img.Set(x, y, color.RGBA64{R: uint16(v), G: uint16(v), B: uint16(v), A: 65535})
			}
		}
	}

	return img
}