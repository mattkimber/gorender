package sprite

import (
	"colour"
	"image"
	"image/color"
	"math"
	"raycaster"
)

type shadeFunc32bpp func(raycaster.RenderSample) (float64, float64, float64)
type shadeFuncIndexed func(raycaster.RenderSample) byte

func ApplyUniformSprite(img *image.RGBA, bounds image.Rectangle, loc image.Point) {
	minX, minY := bounds.Min.X+loc.X, bounds.Min.Y+loc.Y
	maxX, maxY := bounds.Max.X+loc.X, bounds.Max.Y+loc.Y

	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			img.Set(x, y, color.Black)
		}
	}
}

func apply32bppImage(img *image.RGBA, bounds image.Rectangle, loc image.Point, shader shadeFunc32bpp, info raycaster.RenderOutput, softenEdges bool) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			c := get32bppSample(info[x][y], shader, softenEdges)
			img.Set(x+loc.X, y+loc.Y, c)
		}
	}
}

func get32bppSample(info raycaster.RenderInfo, shader shadeFunc32bpp, softenEdges bool) color.RGBA64 {
	total, filled := 0, 0
	cr, cg, cb := 0.0, 0.0, 0.0

	for _, s := range info {
		total++

		if s.Collision {
			r, g, b := shader(s)
			cr += r * r // Summing squares of colours gives more accurate results
			cg += g * g
			cb += b * b
			filled++
		}
	}

	// No collisions = transparent
	if filled == 0 {
		return color.RGBA64{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}
	}

	// Soften edges means that when only some rays collided (typically near edges
	// of an object) we fade to transparent. Otherwise objects are hard-edged, which
	// makes them more likely to suffer aliasing artifacts but also clearer at small
	// sizes
	alpha := 65535
	divisor := float64(filled)
	if softenEdges {
		alpha = (filled * 65535) / (total)
		divisor = float64(total)
	}

	// Return the average colour value
	return color.RGBA64{
		R: clamp(uint16(math.Sqrt(cr / divisor))),
		G: clamp(uint16(math.Sqrt(cg / divisor))),
		B: clamp(uint16(math.Sqrt(cb / divisor))),
		A: uint16(alpha),
	}
}

func clamp(input uint16) uint16 {
	if input < 256 {
		return 256
	}

	return input
}

func applyIndexedImage(img *image.Paletted, pal colour.Palette, bounds image.Rectangle, loc image.Point, shader shadeFuncIndexed, info raycaster.RenderOutput) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			img.SetColorIndex(x+loc.X, y+loc.Y, get8bppSample(info[x][y], shader, pal))
		}
	}
}

func get8bppSample(info raycaster.RenderInfo, shader shadeFuncIndexed, pal colour.Palette) byte {
	values := map[byte]int{}
	specials := 0
	threshold := len(info) / 3

	for _, s := range info {
		if s.Collision {
			idx := shader(s)

			if pal.IsSpecialColour(idx) {
				specials++
				if specials >= threshold {
					return idx
				}
			}

			if idx != 0 {
				values[idx]++
			}
		}
	}

	// TODO: something better than returning the modal index, which produces heavily aliased results
	max := 0
	modalIndex := byte(0)

	for k, v := range values {
		if v > max {
			max = v
			modalIndex = k
		}
	}

	return modalIndex
}
