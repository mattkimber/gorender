package sprite

import (
	"colour"
	"manifest"
	"raycaster"
)

type ShaderInfo struct {
	Colour         colour.RGB
	SpecialColour  colour.RGB
	Alpha          float64
	Specialness    float64
	Normal         colour.RGB
	AveragedNormal colour.RGB
	Depth          colour.RGB
	Occlusion      colour.RGB
	Lighting       colour.RGB
	Shadowing      colour.RGB
	ModalIndex     byte
}

type ShaderOutput [][]ShaderInfo

func GetColour(s *ShaderInfo) colour.RGB {
	return s.Colour
}

func GetNormal(s *ShaderInfo) colour.RGB {
	return s.Normal
}

func GetAveragedNormal(s *ShaderInfo) colour.RGB {
	return s.AveragedNormal
}

func GetDepth(s *ShaderInfo) colour.RGB {
	return s.Depth
}

func GetOcclusion(s *ShaderInfo) colour.RGB {
	return s.Occlusion
}

func GetLighting(s *ShaderInfo) colour.RGB {
	return s.Lighting
}

func GetShadowing(s *ShaderInfo) colour.RGB {
	return s.Shadowing
}

func GetIndex(s *ShaderInfo) byte {
	return s.ModalIndex
}

func GetMaskIndex(s *ShaderInfo) byte {
	if s.Specialness >= 0.5 {
		return s.ModalIndex
	}
	return 0
}

func GetShaderOutput(renderOutput raycaster.RenderOutput, def manifest.Definition, width int, height int) (output ShaderOutput) {
	output = make([][]ShaderInfo, width)

	for x := 0; x < width; x++ {
		output[x] = make([]ShaderInfo, height)

		for y := 0; y < height; y++ {
			output[x][y] = shade(renderOutput[x][y], def)
		}
	}

	return
}

func shade(info raycaster.RenderInfo, def manifest.Definition) (output ShaderInfo) {
	total, filled := 0, 0
	values := map[byte]int{}

	for _, s := range info {
		total++

		if s.Collision {
			output.Colour.Add(Colour(s, def, true))
			output.SpecialColour.Add(Colour(s, def, false))

			if def.Palette.IsSpecialColour(s.Index) {
				output.Specialness += 1.0
			}

			// TODO: in future we will only need this for "special" colours
			if s.Index != 0 {
				values[s.Index]++
			}

			filled++

			if def.Debug {
				output.Normal.Add(Normal(s))
				output.AveragedNormal.Add(AveragedNormal(s))
				output.Depth.Add(Depth(s))
				output.Occlusion.Add(Occlusion(s))
				output.Shadowing.Add(Shadow(s))
				output.Lighting.Add(Lighting(s))
			}
		}
	}

	max := 0

	for k, v := range values {
		if v > max {
			max = v
			output.ModalIndex = k
		}
	}

	// No collisions = transparent
	if filled == 0 {
		return
	}

	// Soften edges means that when only some rays collided (typically near edges
	// of an object) we fade to transparent. Otherwise objects are hard-edged, which
	// makes them more likely to suffer aliasing artifacts but also clearer at small
	// sizes
	output.Alpha = 1.0
	divisor := float64(filled)
	if def.SoftenEdges() {
		output.Alpha = divisor / float64(total)
		divisor = float64(total)
	}

	output.Colour.DivideAndClamp(divisor)
	output.SpecialColour.DivideAndClamp(divisor)

	output.Specialness = output.Specialness / divisor

	if def.Debug {
		output.Normal.DivideAndClamp(divisor)
		output.AveragedNormal.DivideAndClamp(divisor)
		output.Depth.DivideAndClamp(divisor)
		output.Occlusion.DivideAndClamp(divisor)
		output.Shadowing.DivideAndClamp(divisor)
		output.Lighting.DivideAndClamp(divisor)
	}

	return
}
