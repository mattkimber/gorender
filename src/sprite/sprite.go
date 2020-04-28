package sprite

import (
	"image"
	"image/color"
	"raycaster"
	"utils/imageutils"
	"voxelobject"
)

func GetUniformSprite(bounds image.Rectangle) image.Image {
	return imageutils.GetUniformImage(bounds, color.Black)
}

func GetRaycastSprite(object voxelobject.RawVoxelObject, bounds image.Rectangle, angle int) image.Image {
	info := raycaster.GetRaycastOutput(object, angle, bounds.Max.X, bounds.Max.Y)
	img := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if info[x][y].Collision {
				img.Set(x, y, color.Black)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}

	return img
}
