package sprite

import (
	"colour"
	"image"
	"image/color"
	"raycaster"
	"utils/imageutils"
	"voxelobject"
)

func GetUniformSprite(bounds image.Rectangle) image.Image {
	return imageutils.GetUniformImage(bounds, color.Black)
}

func GetRaycastSprite(object voxelobject.RawVoxelObject, pal colour.Palette, bounds image.Rectangle, angle int) image.Image {
	info := raycaster.GetRaycastOutput(object, angle, bounds.Max.X, bounds.Max.Y)
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
