package sprite

import (
	"github.com/mattkimber/gorender/internal/colour"
	"github.com/mattkimber/gorender/internal/manifest"
	"github.com/mattkimber/gorender/internal/raycaster"
	"math"
	"sort"
)

type ShaderInfo struct {
	Colour           colour.RGB
	SpecialColour    colour.RGB
	Alpha            float64
	Specialness      float64
	Normal           colour.RGB
	AveragedNormal   colour.RGB
	Depth            colour.RGB
	Occlusion        colour.RGB
	Lighting         colour.RGB
	Shadowing        colour.RGB
	Detail           colour.RGB
	Transparency     colour.RGB
	Region           int
	LightingCalcDone bool
	DitherChecked    bool
	DitherDone       bool
	MaxLighting      float64
	MinLighting      float64
	ModalIndex       byte
	DitheredIndex    byte
	IsMaskColour     bool
	IsAnimated       bool
	IsBottom         bool
	IsLeft           bool
}

type ShaderOutput [][]ShaderInfo

type RegionInfo struct {
	MinDistanceFromMidpoint float64
	MaxDistanceFromMidpoint float64
	MinIndex                byte
	MaxIndex                byte
	DistanceHistogram       []int
	IndexHistogram          []int
	FilledHistogramBuckets  int
	RangeLength             float64
	Size                    int
	SizeInRange             int
	Range                   *colour.PaletteRange
	RequiresExpansion       bool
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
		R: float64(s.Region % 4 * (65535 / 4)),
		G: float64((s.Region / 4) % 4 * (65535 / 4)),
		B: float64((s.Region / 16) % 4 * (65535 / 4)),
	}
}

func GetShaderOutput(renderOutput raycaster.RenderOutput, spr manifest.Sprite, def *manifest.Definition, width int, height int) (output ShaderOutput) {
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

			identifyRegions(&output, def, currentRegion, x, y, width, height, output[x][y].ModalIndex, &def.Palette, paletteRange)

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

			bestIndex := ditherOutput(def, output, x, y, errCurr, primaryCCPalette, secondaryCCPalette, regularPalette, errNext)

			// Update the range stats
			ditheredRange := def.Palette.Entries[bestIndex].Range

			// Hard reset non-renderable colours for transparent sections
			if bestIndex == 0 {
				output[x][y].ModalIndex = 0
			}

			info := regions[output[x][y].Region]
			info.Size++

			if ditheredRange == info.Range && bestIndex != 0 {
				if info.DistanceHistogram == nil {
					info.DistanceHistogram = make([]int, 256)
				}

				if info.IndexHistogram == nil {
					info.IndexHistogram = make([]int, 256)
				}

				info.SizeInRange++
				info.IndexHistogram[bestIndex]++

				if bestIndex < info.MinIndex || info.MinIndex == 0 {
					info.MinIndex = bestIndex
				}

				if bestIndex > info.MaxIndex {
					info.MaxIndex = bestIndex
				}

				regions[output[x][y].Region] = info
			}
		}

		// Swap the next and current error lines
		errCurr, errNext = errNext, errCurr
	}

	for idx, region := range regions {

		if region.SizeInRange > 1 {
			rng := region.Range
			rangeLength := float64(rng.End - rng.Start)
			region.RangeLength = rangeLength

			// Check the histogram. If more than 50% of a 4px or larger region is a single colour, we need to expand
			// the colour range
			//
			// Additionally if the region is >8px and only two colours have been used, we will expand the range
			histogramMax := 0
			usedColours := 0

			for i := range region.IndexHistogram {
				if region.IndexHistogram[i] > 0 {
					usedColours++
				}

				if region.IndexHistogram[i] > histogramMax {
					histogramMax = region.IndexHistogram[i]
				}
			}

			// Expand the range if not many colours are used
			if region.MaxIndex-region.MinIndex < rng.ExpectedColourRange {
				if region.MinIndex > rng.Start {
					region.MinIndex = region.MinIndex - 1
				}

				if region.MaxIndex < rng.End {
					region.MaxIndex = region.MaxIndex + 1
				}
			}

			regions[idx] = region
		}
	}

	// Do the second pass dithered output to expand the colour range
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			paletteRange := def.Palette.Entries[output[x][y].DitheredIndex].Range
			if paletteRange == nil || paletteRange.IsAnimatedLight || paletteRange.IsNonRenderable {
				// Don't alter special colours
				continue
			}

			if def.Manifest.DitherFlatAreas {
				minLighting, maxLighting, totalPixels := math.MaxFloat64, 0.0, 0
				getLightingForSameColourArea(&output, def, x, y, width, height, output[x][y].ModalIndex, &minLighting, &maxLighting, &totalPixels)

				if maxLighting-minLighting > 0.0 {
					// Work out how much to dither so 60% of output is un-dithered, 20% darkened, 20% lightened
					lightingValues := make([]float64, 0)
					getLightingValues(&output, def, x, y, width, height, output[x][y].ModalIndex, minLighting, maxLighting, &lightingValues)

					if len(lightingValues) > 1 {
						sort.Slice(lightingValues, func(i, j int) bool {
							return lightingValues[i] < lightingValues[j]
						})

						ditherThresholdLow := lightingValues[len(lightingValues)/5]
						ditherThresholdHigh := lightingValues[(len(lightingValues)*4)/5]

						if ditherThresholdLow != ditherThresholdHigh {
							doColourPush(&output, def, x, y, width, height, output[x][y].ModalIndex, &def.Palette, minLighting, maxLighting, ditherThresholdLow, ditherThresholdHigh)
						}
					}
				}
			}

			// "Fosterise" by darkening pixels at the bottom and left
			if def.Manifest.Fosterise && (output[x][y].IsBottom || output[x][y].IsLeft) && output[x][y].DitheredIndex > paletteRange.Start {
				output[x][y].DitheredIndex--
			}
		}
	}

	return
}

func ditherOutput(def *manifest.Definition, output ShaderOutput, x int, y int, errCurr []colour.RGB, primaryCCPalette []colour.RGB, secondaryCCPalette []colour.RGB, regularPalette []colour.RGB, errNext []colour.RGB) (bestIndex byte) {
	var ditherError colour.RGB

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

func floodFill(x, y, width, height int, fn func(int, int)) {
	// Recursively flood fill in the adjacent directions
	if x > 0 {
		fn(x-1, y)
	}

	if y > 0 {
		fn(x, y-1)
	}

	if x < width-1 {
		fn(x+1, y)
	}

	if y < height-1 {
		fn(x, y+1)
	}

}

func identifyRegions(output *ShaderOutput, def *manifest.Definition, region int, x, y int, width, height int, previousIndex byte, palette *colour.Palette, paletteRange *colour.PaletteRange) {
	index := (*output)[x][y].ModalIndex
	thisRegion := (*output)[x][y].Region
	thisRange := (*palette).Entries[index].Range

	gap := int(previousIndex) - int(index)
	if gap < 0 {
		gap = -gap
	}

	// If not the same palette range, or we already set the region, or the indexes are too far apart, return
	if thisRange != paletteRange || thisRegion == region || gap > paletteRange.MaxGapInRegion {
		return
	}

	(*output)[x][y].Region = region

	// Recursively find the regions
	floodFill(x, y, width, height, func(x1, y1 int) {
		identifyRegions(output, def, region, x1, y1, width, height, index, palette, paletteRange)
	})

	if x > 0 && x < width-1 && (*output)[x-1][y].Region != region && (*output)[x+1][y].Region == region {
		if !def.Manifest.NoEdgeFosterisation || (*output)[x-1][y].ModalIndex != 0 {
			(*output)[x][y].IsLeft = true
		}
	}

	// Left edge of sprite at border
	if x == 0 && x < width-1 && (*output)[x+1][y].Region == region {
		if !def.Manifest.NoEdgeFosterisation {
			(*output)[x][y].IsLeft = true
		}
	}

	if y > 0 && y < height-1 && (*output)[x][y+1].Region != region && (*output)[x][y-1].Region == region {
		if !def.Manifest.NoEdgeFosterisation || (*output)[x][y+1].ModalIndex != 0 {
			(*output)[x][y].IsBottom = true
		}
	}

	// Bottom edge of sprite at border
	if y == height-1 && y > 0 && (*output)[x][y-1].Region == region {
		if !def.Manifest.NoEdgeFosterisation {
			(*output)[x][y].IsBottom = true
		}
	}

	return
}

func getLightingForSameColourArea(output *ShaderOutput, def *manifest.Definition, x, y int, width, height int, previousIndex byte, minLighting, maxLighting *float64, totalPixels *int) {
	index := (*output)[x][y].DitheredIndex

	if (*output)[x][y].LightingCalcDone || index != previousIndex {
		return
	}

	if (*output)[x][y].Lighting.R > *maxLighting {
		*maxLighting = (*output)[x][y].Lighting.R
	}

	if (*output)[x][y].Lighting.R < *minLighting {
		*minLighting = (*output)[x][y].Lighting.R
	}

	(*output)[x][y].LightingCalcDone = true
	*totalPixels++

	// Recursively flood fill in the adjacent directions
	floodFill(x, y, width, height, func(x1, y1 int) {
		getLightingForSameColourArea(output, def, x1, y1, width, height, index, minLighting, maxLighting, totalPixels)
	})

	return
}

func getLightingValues(output *ShaderOutput, def *manifest.Definition, x, y int, width, height int, previousIndex byte, minLighting, maxLighting float64, lightingValues *[]float64) {
	index := (*output)[x][y].DitheredIndex

	if (*output)[x][y].DitherChecked || index != previousIndex {
		return
	}

	lightingValue := ((*output)[x][y].Lighting.R - minLighting) / (maxLighting - minLighting)
	*lightingValues = append(*lightingValues, lightingValue)
	(*output)[x][y].DitherChecked = true

	// Recursively flood fill in the adjacent directions
	floodFill(x, y, width, height, func(x1, y1 int) {
		getLightingValues(output, def, x1, y1, width, height, index, minLighting, maxLighting, lightingValues)
	})

	return
}

func doColourPush(output *ShaderOutput, def *manifest.Definition, x, y int, width, height int, previousIndex byte, palette *colour.Palette, minLighting, maxLighting, ditherThresholdLow, ditherThresholdHigh float64) {
	index := (*output)[x][y].DitheredIndex
	thisRange := (*palette).Entries[index].Range

	if (*output)[x][y].DitherDone || index != previousIndex {
		return
	}

	lightingValue := ((*output)[x][y].Lighting.R - minLighting) / (maxLighting - minLighting)
	if lightingValue > ditherThresholdHigh && (x%2+y)%2 == 0 && index < thisRange.End {
		(*output)[x][y].DitheredIndex++
	} else if lightingValue < ditherThresholdLow && (x%2+y)%2 == 0 && index > thisRange.Start {
		(*output)[x][y].DitheredIndex--
	}

	(*output)[x][y].DitherDone = true

	// Recursively flood fill in the adjacent directions
	floodFill(x, y, width, height, func(x1, y1 int) {
		doColourPush(output, def, x1, y1, width, height, index, palette, minLighting, maxLighting, ditherThresholdLow, ditherThresholdHigh)
	})

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

func shade(info raycaster.RenderInfo, def *manifest.Definition, prevIndex byte) (output ShaderInfo) {
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

			output.Lighting = output.Lighting.Add(Lighting(s).MultiplyBy(s.Influence))

			if def.Debug {
				floatCount := float64(s.Count)
				output.Normal = output.Normal.Add(Normal(s).MultiplyBy(floatCount))
				output.AveragedNormal = output.AveragedNormal.Add(AveragedNormal(s).MultiplyBy(floatCount))
				output.Depth = output.Depth.Add(Depth(s).MultiplyBy(floatCount))
				output.Occlusion = output.Occlusion.Add(Occlusion(s).MultiplyBy(floatCount))
				output.Shadowing = output.Shadowing.Add(Shadow(s).MultiplyBy(floatCount))
				output.Detail = output.Detail.Add(Detail(s).MultiplyBy(floatCount))
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
	if totalSamples == 0 || filledSamples*100/totalSamples <= hardEdgeThreshold {
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

	output.Lighting.DivideAndClamp(divisor)

	if def.Debug {
		debugDivisor := float64(filledSamples)
		output.Normal.DivideAndClamp(debugDivisor)
		output.AveragedNormal.DivideAndClamp(debugDivisor)
		output.Depth.DivideAndClamp(debugDivisor)
		output.Occlusion.DivideAndClamp(debugDivisor)
		output.Shadowing.DivideAndClamp(debugDivisor)
		output.Detail.DivideAndClamp(debugDivisor)
		output.Transparency = FloatValue(float64(filledSamples) / float64(totalSamples))
	}

	return
}
