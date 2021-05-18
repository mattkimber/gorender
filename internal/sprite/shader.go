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
	Region         int
	ModalIndex     byte
	DitheredIndex  byte
	IsMaskColour   bool
	IsAnimated     bool
	Midpoint 	   float64
	DistanceFromMidpoint float64
}

type ShaderOutput [][]ShaderInfo

type RegionInfo struct {
	MinDistanceFromMidpoint float64
	MaxDistanceFromMidpoint float64
	LowColourPush float64
	HighColourPush float64
	RangeLength float64
	Size int
	SizeInRange int
	Range *colour.PaletteRange
}

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

func GetRegion(s *ShaderInfo) colour.RGB {
	return colour.RGB{
		R: float64(s.Region % 4 * (65535/4)),
		G: float64((s.Region/4) % 4 * (65535/4)),
		B: float64((s.Region/16) % 4 * (65535/4)),
	}
}

func GetMidpointDistance(s *ShaderInfo) colour.RGB {
	return FloatValue(s.DistanceFromMidpoint)
}

func GetShaderOutput(renderOutput raycaster.RenderOutput, spr manifest.Sprite, def manifest.Definition, width int, height int) (output ShaderOutput) {
	output = make([][]ShaderInfo, width)

	xoffset, yoffset := int(spr.OffsetX*def.Scale), int(spr.OffsetY*def.Scale)


	prevIndex := byte(0)

	for x := 0; x < width; x++ {
		output[x] = make([]ShaderInfo, height)

		for y := 0; y < height; y++ {
			rx := x + xoffset
			ry := y + yoffset
			if rx < 0 || rx >= width || ry < 0 || ry >= height {
				continue
			}

			if x > 1 {
				prevIndex = output[x-1][y].ModalIndex
			} else {
				prevIndex = 0
			}

			output[x][y] = shade(renderOutput[rx][ry], def, prevIndex)

		}
	}

	currentRegion := 1
	regions := make(map[int]RegionInfo)

	// Calculate regions from the shaded output
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			info := RegionInfo{}

			// No region for transparent/empty voxels
			if output[x][y].ModalIndex == 0 {
				continue
			}

			// Don't set region if it was already set
			if output[x][y].Region != 0 {
				continue
			}

			// Flood fill the region connected to this pixel
			paletteRange := def.Palette.Entries[output[x][y].ModalIndex].Range
			info.Range = paletteRange

			floodFill(&output, currentRegion, x, y, width, height, &def.Palette, paletteRange)

			regions[currentRegion] = info
			currentRegion++
		}
	}


	// Floyd-Steinberg error rows
	errCurr := make([]colour.RGB, height+2)
	errNext := make([]colour.RGB, height+2)

	// Palettes
	regularPalette := def.Palette.GetRegularPalette()
	primaryCCPalette := def.Palette.GetPrimaryCompanyColourPalette()
	secondaryCCPalette := def.Palette.GetSecondaryCompanyColourPalette()

	// Get the first pass dithered output to find what the colour ranges are
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			bestIndex, ditherErr := ditherOutput(def, output, x, y, errCurr, primaryCCPalette, secondaryCCPalette, regularPalette, errNext)

			// Update the range stats
			ditheredRange := def.Palette.Entries[bestIndex].Range

			if ditheredRange != nil {
				position := 0.0

				if bestIndex != ditheredRange.Start && bestIndex != ditheredRange.End {
					otherIndex := bestIndex + 1

					resultErr := ditherErr.Subtract(def.Palette.Entries[bestIndex].GetRGB())
					otherResultErr := ditherErr.Subtract(def.Palette.Entries[otherIndex].GetRGB())

					distance := math.Sqrt(resultErr.R*resultErr.R) + math.Sqrt(resultErr.G*resultErr.G) + math.Sqrt(resultErr.B*resultErr.B)
					otherDistance := math.Sqrt(otherResultErr.R*otherResultErr.R) + math.Sqrt(otherResultErr.G*otherResultErr.G) + math.Sqrt(otherResultErr.B*otherResultErr.B)

					position = distance / (distance + otherDistance)
				}

				rangeLength := float64(ditheredRange.End - ditheredRange.Start)
				midpoint := float64(ditheredRange.Start) + (rangeLength / 2.0)
				output[x][y].Midpoint = midpoint
				output[x][y].DistanceFromMidpoint = ((float64(bestIndex) + position) - midpoint) / (rangeLength / 2.0)
			}

			info := regions[output[x][y].Region]
			info.Size++

			if ditheredRange == info.Range && bestIndex != 0 {

				info.SizeInRange++

				if output[x][y].DistanceFromMidpoint < info.MinDistanceFromMidpoint {
					info.MinDistanceFromMidpoint = output[x][y].DistanceFromMidpoint
				}

				if output[x][y].DistanceFromMidpoint > info.MaxDistanceFromMidpoint {
					info.MaxDistanceFromMidpoint = output[x][y].DistanceFromMidpoint
				}

				regions[output[x][y].Region] = info
			}
		}

		// Swap the next and current error lines
		errCurr, errNext = errNext, errCurr
	}

	for idx, region := range regions {
		if region.Size > 1 {
			rng := region.Range
			rangeLength := float64(rng.End - rng.Start)

			region.RangeLength = rangeLength

			// TODO: should be configurable in manifest
			const MAX_PUSH_AMOUNT = 2.0

			region.HighColourPush = 1.0
			region.LowColourPush = 1.0

			if region.MaxDistanceFromMidpoint < 1.0 && region.MaxDistanceFromMidpoint > 0.1 {
				region.HighColourPush = 1.0 / region.MaxDistanceFromMidpoint
				if region.HighColourPush > MAX_PUSH_AMOUNT {
					region.HighColourPush = MAX_PUSH_AMOUNT
				}
			}

			if region.MinDistanceFromMidpoint > -1.0 && region.MinDistanceFromMidpoint < -0.1 {
				region.LowColourPush = -1.0 / region.MinDistanceFromMidpoint
				if region.LowColourPush > MAX_PUSH_AMOUNT {
					region.LowColourPush = MAX_PUSH_AMOUNT
				}
			}

			regions[idx] = region
		}
	}


	// Do the second pass dithered output to expand the colour range
	// TODO: do something useful
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			region, ok := regions[output[x][y].Region]
			paletteRange := def.Palette.Entries[output[x][y].DitheredIndex].Range

			if ok && region.SizeInRange > 1 && output[x][y].DitheredIndex != 0 && region.RangeLength > 0 &&
				output[x][y].Midpoint != 0 && region.Range == paletteRange {

				distance := output[x][y].DistanceFromMidpoint * (region.RangeLength/2.0)
				if float64(output[x][y].DitheredIndex) > output[x][y].Midpoint {
					distance = distance * region.HighColourPush
				} else {
					distance = distance * region.LowColourPush
				}

				output[x][y].DitheredIndex = byte(math.Round(output[x][y].Midpoint + distance))
				if output[x][y].DitheredIndex < region.Range.Start {
					output[x][y].DitheredIndex = region.Range.Start
				} else if output[x][y].DitheredIndex > region.Range.End {
					output[x][y].DitheredIndex = region.Range.End
				}

			}
		}
	}

	return
}

func ditherOutput(def manifest.Definition, output ShaderOutput, x int, y int, errCurr []colour.RGB, primaryCCPalette []colour.RGB, secondaryCCPalette []colour.RGB, regularPalette []colour.RGB, errNext []colour.RGB) (bestIndex byte, ditherError colour.RGB) {
	rng := def.Palette.Entries[output[x][y].ModalIndex].Range
	if rng == nil {
		rng = &colour.PaletteRange{}
	}

	if output[x][y].Alpha < def.Manifest.EdgeThreshold {
		bestIndex = 0
	} else if rng.IsPrimaryCompanyColour {
		if y > 0 && def.Palette.IsSpecialColour(output[x][y-1].ModalIndex) {
			ditherError = output[x][y].SpecialColour
		} else {
			ditherError = output[x][y].SpecialColour.Add(errCurr[y+1])
		}
		bestIndex = getBestIndex(ditherError, primaryCCPalette)
	} else if rng.IsSecondaryCompanyColour {
		if y > 0 && def.Palette.IsSpecialColour(output[x][y-1].ModalIndex) {
			ditherError = output[x][y].SpecialColour
		} else {
			ditherError = output[x][y].SpecialColour.Add(errCurr[y+1])
		}
		bestIndex = getBestIndex(ditherError, secondaryCCPalette)
	} else if rng.IsAnimatedLight {
		output[x][y].IsAnimated = true
		// Never add error values to special colours
		bestIndex = output[x][y].ModalIndex
		ditherError = def.Palette.Entries[bestIndex].GetRGB()
	} else {
		if y > 0 && def.Palette.IsSpecialColour(output[x][y-1].ModalIndex) {
			ditherError = output[x][y].Colour
		} else {
			ditherError = output[x][y].Colour.Add(errCurr[y+1])
		}
		bestIndex = getBestIndex(ditherError, regularPalette)
	}

	output[x][y].DitheredIndex = bestIndex

	if def.Palette.IsSpecialColour(bestIndex) {
		output[x][y].IsMaskColour = true
	}

	resultError := colour.RGB{}

	if output[x][y].Alpha >= def.Manifest.EdgeThreshold {
		resultError = colour.PermissiveClampRGB(ditherError.Subtract(def.Palette.Entries[bestIndex].GetRGB()))
	}

	// Apply Floyd-Steinberg error
	errNext[y+0] = errNext[y+0].Add(resultError.MultiplyBy(3.0 / 16))
	errNext[y+1] = errNext[y+1].Add(resultError.MultiplyBy(5.0 / 16))
	errNext[y+2] = errNext[y+2].Add(resultError.MultiplyBy(1.0 / 16))
	errCurr[y+2] = errCurr[y+2].Add(resultError.MultiplyBy(7.0 / 16))

	errCurr[y+1] = colour.RGB{}
	return
}

func floodFill(output *ShaderOutput, region int, x, y int, width, height int, palette *colour.Palette, paletteRange *colour.PaletteRange) bool {
	index := (*output)[x][y].ModalIndex
	thisRegion := (*output)[x][y].Region
	thisRange := (*palette).Entries[index].Range

	// If not the same palette range, or we already set the region, return
	if thisRange != paletteRange || thisRegion == region {
		return false
	}

	(*output)[x][y].Region = region

	// Recursively flood fill in the adjacent directions
	if x > 0 {
		floodFill(output, region, x - 1, y, width, height, palette, paletteRange)
	}

	if y > 0 {
		floodFill(output, region, x, y - 1, width, height, palette, paletteRange)
	}

	if x < width - 1 {
		floodFill(output, region, x + 1, y, width, height, palette, paletteRange)
	}

	if y < height - 1 {
		floodFill(output, region, x, y + 1, width, height, palette, paletteRange)
	}

	return true
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

func shade(info raycaster.RenderInfo, def manifest.Definition, prevIndex byte) (output ShaderInfo) {
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
			filledSamples += s.Count

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
				// Loop makes this a little slower but is fine for debug purposes
				for i := 0; i < s.Count; i++ {
					output.Normal = output.Normal.Add(Normal(s))
					output.AveragedNormal = output.AveragedNormal.Add(AveragedNormal(s))
					output.Depth = output.Depth.Add(Depth(s))
					output.Occlusion = output.Occlusion.Add(Occlusion(s))
					output.Shadowing = output.Shadowing.Add(Shadow(s))
					output.Lighting = output.Lighting.Add(Lighting(s))
					output.Detail = output.Detail.Add(Detail(s))
				}
			}
		}

		totalSamples = totalSamples + s.Count
	}

	max := 0.0
	alternateModal := byte(0)

	for k, v := range values {
		if v > max {
			max = v
			// Store the previous modal
			alternateModal = output.ModalIndex
			output.ModalIndex = k
		}
	}

	// Supply a same-range alternative if we are going to repeat the same colour and we have an alternative
	if output.ModalIndex == prevIndex && def.Palette.Entries[output.ModalIndex].Range == def.Palette.Entries[alternateModal].Range && alternateModal != 0 {
		output.ModalIndex = alternateModal
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
