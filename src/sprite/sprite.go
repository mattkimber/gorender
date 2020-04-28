package sprite

import (
	"image"
	"image/color"
	"utils/imageutils"
)

func GetSprite(bounds image.Rectangle, angle int) image.Image {
	return imageutils.GetUniformImage(bounds, color.Black)
}
