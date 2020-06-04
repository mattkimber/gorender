package colour

import "image/color"

type RGB struct {
	R float64
	G float64
	B float64
}

func (rgb *RGB) DivideAndClamp(divisor float64) {
	rgb.R = Clamp(rgb.R / divisor)
	rgb.G = Clamp(rgb.G / divisor)
	rgb.B = Clamp(rgb.B / divisor)
}

func (rgb *RGB) GetRGBA(alpha float64) color.RGBA64 {
	return color.RGBA64{
		R: uint16(rgb.R),
		G: uint16(rgb.G),
		B: uint16(rgb.B),
		A: uint16(alpha * 65535),
	}
}

func (rgb RGB) Add(input RGB) (result RGB) {
	result.R = rgb.R + input.R
	result.G = rgb.G + input.G
	result.B = rgb.B + input.B

	return
}

func (rgb RGB) Subtract(input RGB) (result RGB) {
	result.R = rgb.R - input.R
	result.G = rgb.G - input.G
	result.B = rgb.B - input.B

	return
}

func (rgb RGB) MultiplyBy(value float64) (result RGB) {
	result.R = rgb.R * value
	result.G = rgb.G * value
	result.B = rgb.B * value

	return
}

func FromPaletteEntry(p PaletteEntry) RGB {
	return RGB{
		R: float64(p.R) * 255,
		G: float64(p.G) * 255,
		B: float64(p.B) * 255,
	}
}
