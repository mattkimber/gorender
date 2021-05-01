package sprite

import (
	"github.com/mattkimber/gorender/internal/colour"
	"github.com/mattkimber/gorender/internal/manifest"
	"github.com/mattkimber/gorender/internal/raycaster"
	"math"
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
	Detail         colour.RGB
	Transparency   colour.RGB
	ModalIndex     byte
	DitheredIndex  byte
	IsMaskColour   bool
	IsAnimated     bool
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

func GetDetail(s *ShaderInfo) colour.RGB {
	return s.Detail
}

func GetTransparency(s *ShaderInfo) colour.RGB {
	return s.Transparency
}

func GetIndex(s *ShaderInfo) byte {
	return s.DitheredIndex
}

func GetMaskIndex(s *ShaderInfo) byte {
	if s.Specialness > 0.75 || s.IsAnimated {
		return s.ModalIndex
	} else if s.Specialness > 0.25 && s.IsMaskColour {
		return s.DitheredIndex
	}
	return 0
}

func GetShaderOutput(renderOutput raycaster.RenderOutput, spr manifest.Sprite, def manifest.Definition, width int, height int) (output ShaderOutput) {
	output = make([][]ShaderInfo, width)

	xoffset, yoffset := int(spr.OffsetX*def.Scale), int(spr.OffsetY*def.Scale)

	// Palettes
	regularPalette := def.Palette.GetRegularPalette()
	primaryCCPalette := def.Palette.GetPrimaryCompanyColourPalette()
	secondaryCCPalette := def.Palette.GetSecondaryCompanyColourPalette()

	// Floyd-Steinberg error rows
	errCurr := make([]colour.RGB, height+2)
	errNext := make([]colour.RGB, height+2)

	var error colour.RGB

	for x := 0; x < width; x++ {
		output[x] = make([]ShaderInfo, height)

		for y := 0; y < height; y++ {
			rx := x + xoffset
			ry := y + yoffset
			if rx < 0 || rx >= width || ry < 0 || ry >= height {
				continue
			}

			output[x][y] = shade(renderOutput[rx][ry], def)
			bestIndex := byte(0)

			rng := def.Palette.Entries[output[x][y].ModalIndex].Range
			if rng == nil {
				rng = &colour.PaletteRange{}
			}

			if output[x][y].Alpha < def.Manifest.EdgeThreshold {
				bestIndex = 0
			} else if rng.IsPrimaryCompanyColour {
				if y > 0 && def.Palette.IsSpecialColour(output[x][y-1].ModalIndex) {
					error = output[x][y].SpecialColour
				} else {
					error = output[x][y].SpecialColour.Add(errCurr[y+1])
				}
				bestIndex = getBestIndex(error, primaryCCPalette)
			} else if rng.IsSecondaryCompanyColour {
				if y > 0 && def.Palette.IsSpecialColour(output[x][y-1].ModalIndex) {
					error = output[x][y].SpecialColour
				} else {
					error = output[x][y].SpecialColour.Add(errCurr[y+1])
				}
				bestIndex = getBestIndex(error, secondaryCCPalette)
			} else if rng.IsAnimatedLight {
				output[x][y].IsAnimated = true
				// Never add error values to special colours
				bestIndex = output[x][y].ModalIndex
				error = def.Palette.Entries[bestIndex].GetRGB()
			} else {
				if y > 0 && def.Palette.IsSpecialColour(output[x][y-1].ModalIndex) {
					error = output[x][y].Colour
				} else {
					error = output[x][y].Colour.Add(errCurr[y+1])
				}
				bestIndex = getBestIndex(error, regularPalette)
			}

			output[x][y].DitheredIndex = bestIndex

			if def.Palette.IsSpecialColour(bestIndex) {
				output[x][y].IsMaskColour = true
			}

			if output[x][y].Alpha >= def.Manifest.EdgeThreshold {
				error = colour.ClampRGB(error.Subtract(def.Palette.Entries[bestIndex].GetRGB()))
			} else {
				error = colour.RGB{}
			}

			// Apply Floyd-Steinberg error
			errNext[y+0] = errNext[y+0].Add(error.MultiplyBy(3.0 / 16))
			errNext[y+1] = errNext[y+1].Add(error.MultiplyBy(5.0 / 16))
			errNext[y+2] = errNext[y+2].Add(error.MultiplyBy(1.0 / 16))
			errCurr[y+2] = errCurr[y+2].Add(error.MultiplyBy(7.0 / 16))

			errCurr[y+1] = colour.RGB{}
		}

		// Swap the next and current error lines
		errCurr, errNext = errNext, errCurr
	}

	return
}

func getBestIndex(error colour.RGB, palette []colour.RGB) byte {
	bestIndex, bestSum := 0, math.MaxFloat64
	for index, p := range palette {
		if p.R > 65000 && (p.G == 0 || p.G > 65000) && p.B > 65000 {
			continue
		}

		sum := squareDiff(error.R, p.R) + squareDiff(error.G, p.G) + squareDiff(error.B, p.B)
		if sum < bestSum {
			bestIndex, bestSum = index, sum
			if sum == 0 {
				break
			}
		}
	}

	return byte(bestIndex)
}

func squareDiff(a, b float64) float64 {
	diff := a - b
	return diff * diff
}

func shade(info raycaster.RenderInfo, def manifest.Definition) (output ShaderInfo) {
	totalInfluence, filledInfluence := 0.0, 0.0
	filledSamples, totalSamples := 0, 0
	values := map[byte]float64{}
	fAccuracy := float64(def.Manifest.Accuracy)
	hardEdgeThreshold := int(def.Manifest.HardEdgeThreshold * 100.0)

	minDepth := math.MaxInt64
	for _, s := range info {
		if s.Collision && s.Depth < minDepth {
			minDepth = s.Depth
		}
	}

	for _, s := range info {
		if s.IsRecovered {
			s.Influence = s.Influence * (1.0 - def.Manifest.RecoveredVoxelSuppression)
		}

		// Voxel samples considered to be more representative of fine details can be boosted
		// to make them more likely to appear in the output.
		if def.Manifest.DetailBoost != 0 {
			s.Influence = s.Influence * (1.0 + (s.Detail * def.Manifest.DetailBoost))
		}

		// Boost samples closest to the camera
		if s.Depth != minDepth {
			s.Influence = s.Influence / fAccuracy
		}

		totalInfluence += s.Influence

		if s.Collision && def.Palette.IsRenderable(s.Index) {
			filledInfluence += s.Influence
			filledSamples += 1

			output.Colour = output.Colour.Add(Colour(s, def, true, s.Influence))
			output.SpecialColour = output.SpecialColour.Add(Colour(s, def, false, s.Influence))

			if def.Palette.IsSpecialColour(s.Index) {
				output.Specialness += 1.0 * s.Influence
				values[s.Index]++
			}

			if s.Index != 0 {
				values[s.Index] += s.Influence
			}

			if def.Debug {
				output.Normal = output.Normal.Add(Normal(s))
				output.AveragedNormal = output.AveragedNormal.Add(AveragedNormal(s))
				output.Depth = output.Depth.Add(Depth(s))
				output.Occlusion = output.Occlusion.Add(Occlusion(s))
				output.Shadowing = output.Shadowing.Add(Shadow(s))
				output.Lighting = output.Lighting.Add(Lighting(s))
				output.Detail = output.Detail.Add(Detail(s))
			}
		}

		totalSamples++
	}

	max := 0.0

	for k, v := range values {
		if v > max {
			max = v
			output.ModalIndex = k
		}
	}

	// Fewer than hard edge threshold collisions = transparent
	if totalSamples == 0 || filledSamples * 100 / totalSamples <= hardEdgeThreshold {
		return ShaderInfo{}
	}

	// Soften edges means that when only some rays collided (typically near edges
	// of an object) we fade to transparent. Otherwise objects are hard-edged, which
	// makes them more likely to suffer aliasing artifacts but also clearer at small
	// sizes
	output.Alpha = 1.0
	divisor := filledInfluence

	if def.SoftenEdges() {
		output.Alpha = divisor / totalInfluence
	}

	if def.Manifest.FadeToBlack {
		divisor = totalInfluence
	}

	output.Colour.DivideAndClamp(divisor)
	output.SpecialColour.DivideAndClamp(divisor)

	output.Specialness = output.Specialness / divisor

	if def.Debug {
		debugDivisor := float64(filledSamples)
		output.Normal.DivideAndClamp(debugDivisor)
		output.AveragedNormal.DivideAndClamp(debugDivisor)
		output.Depth.DivideAndClamp(debugDivisor)
		output.Occlusion.DivideAndClamp(debugDivisor)
		output.Shadowing.DivideAndClamp(debugDivisor)
		output.Lighting.DivideAndClamp(debugDivisor)
		output.Detail.DivideAndClamp(debugDivisor)
		output.Transparency = FloatValue(float64(filledSamples)/ float64(totalSamples))
	}

	return
}
