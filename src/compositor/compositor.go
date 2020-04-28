package compositor

import (
	"fmt"
	"image"
	"image/draw"
)

func Composite(src image.Image, dst image.Image, loc image.Point, size image.Rectangle) error {
	writableDst, ok := dst.(draw.Image)
	if !ok {
		return fmt.Errorf("could not convert destination image to writable image")
	}

	rect := image.Rectangle{Min: loc, Max: src.Bounds().Add(loc).Max}
	draw.Draw(writableDst, rect, src, image.Point{}, draw.Src)

	return nil
}
