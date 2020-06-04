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

func (rgb *RGB) Add(input RGB) {
	rgb.R += input.R
	rgb.G += input.G
	rgb.B += input.B
}
