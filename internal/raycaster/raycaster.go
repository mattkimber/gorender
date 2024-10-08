package raycaster

import (
	"github.com/mattkimber/gorender/internal/geometry"
	"github.com/mattkimber/gorender/internal/manifest"
	"github.com/mattkimber/gorender/internal/sampler"
	"github.com/mattkimber/gorender/internal/voxelobject"
	"sync"
)

type RenderInfo []RenderSample

type RenderSample struct {
	Collision              bool
	Index                  byte
	Normal, AveragedNormal geometry.Vector3
	Depth, Occlusion       int
	LightAmount            float64
	Shadowing              float64
	Influence              float64
	Detail                 float64
	Count                  int
	IsRecovered            bool
}

type RayResult struct {
	X, Y, Z               int
	HasGeometry           bool
	Depth                 int
	IsRecovered           bool
	ApproachedBoundingBox bool
}

type RenderOutput [][]RenderInfo

func GetRaycastOutput(object voxelobject.ProcessedVoxelObject, m manifest.Manifest, spr manifest.Sprite, sampler sampler.Samples) RenderOutput {
	size := object.Size

	// Handle slicing functionality
	minX, maxX := 0, object.Size.X
	if m.SliceLength > 0 && m.SliceThreshold > 0 && m.SliceThreshold < object.Size.X {
		midpoint := (object.Size.X / 2) - (m.SliceLength / 2)
		minX = midpoint - (m.SliceLength * spr.Slice)
		maxX = minX + m.SliceLength

		// Allow sprites to overlap to avoid edge transparency effects
		minX -= m.SliceOverlap
		maxX += m.SliceOverlap

		if minX < 0 {
			minX = 0
		}
		if maxX > object.Size.X {
			maxX = object.Size.X
		}
	}

	limits := geometry.Vector3{X: float64(size.X), Y: float64(size.Y), Z: float64(size.Z)}

	viewport := getViewportPlane(spr.Angle, m, spr.ZError, size, float64(spr.RenderElevationAngle))
	ray := geometry.Zero().Subtract(getRenderDirection(spr.Angle, float64(spr.RenderElevationAngle)))

	lighting := getLightingDirection(spr.Angle+float64(m.LightingAngle), float64(m.LightingElevation), spr.Flip)
	result := make(RenderOutput, len(sampler))

	wg := sync.WaitGroup{}
	wg.Add(sampler.Width())

	joggle := spr.Joggle + m.Joggle

	w, h := sampler.Width(), sampler.Height()

	for x := 0; x < w; x++ {
		thisX := x
		go func() {
			result[thisX] = make([]RenderInfo, h)
			for y := 0; y < h; y++ {
				samples := sampler[thisX][y]
				result[thisX][y] = make(RenderInfo, len(samples))
				raycastSamples(viewport, &samples, ray, limits, object, m, spr, lighting, result, thisX, y, minX, maxX, joggle)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return result
}

func raycastSamples(
	viewport geometry.Plane,
	samples *sampler.SampleList,
	ray geometry.Vector3,
	limits geometry.Vector3,
	object voxelobject.ProcessedVoxelObject,
	m manifest.Manifest,
	spr manifest.Sprite,
	lighting geometry.Vector3,
	result RenderOutput,
	thisX int,
	y int,
	minX int,
	maxX int,
	joggle float64) {

	px, py, pz, pi := 0, 0, 0, 0

	for i := range *samples {
		result[thisX][y][i].Count = 1
	}

	for i, s := range *samples {
		loc0 := viewport.BiLerpWithinPlane(s.Location.X, s.Location.Y)
		loc0.Z += joggle
		loc := getIntersectionWithBounds(loc0, ray, limits)

		rayResult := castFpRay(object, loc0, loc, ray, limits, spr.Flip)

		if rayResult.HasGeometry && rayResult.X >= minX && rayResult.X <= maxX {
			// Speed up for cases where we already encountered this voxel - reduce the amount of sampling needed
			// later
			if rayResult.X == px && rayResult.Y == py && rayResult.Z == pz {
				result[thisX][y][pi].Influence += s.Influence
				result[thisX][y][pi].Count++

				// Set the count for this element to 0
				result[thisX][y][i].Count = 0
				continue
			} else {
				px = rayResult.X
				py = rayResult.Y
				pz = rayResult.Z
				pi = i
			}

			shadowResult := 0
			if getLightingValue(object.Elements[rayResult.X][rayResult.Y][rayResult.Z].AveragedNormal, lighting) > m.ShadowThreshold {
				resultVec := geometry.Vector3{X: float64(rayResult.X), Y: float64(rayResult.Y), Z: float64(rayResult.Z)}
				shadowLoc := resultVec

				shadowVec := geometry.Zero().Subtract(lighting).Normalise()

				for {
					sx, sy, sz := int(shadowLoc.X), int(shadowLoc.Y), int(shadowLoc.Z)

					if sx != rayResult.X || sy != rayResult.Y || sz != rayResult.Z {
						break
					}

					shadowLoc = shadowLoc.Add(shadowVec)
				}

				// Don't flip Y when calculating shadows, as it has been pre-flipped on input.
				shadowResult = castFpRay(object, shadowLoc, shadowLoc, shadowVec, limits, false).Depth
			}
			setResult(&result[thisX][y][i], object.Elements[rayResult.X][rayResult.Y][rayResult.Z], lighting, rayResult.Depth, shadowResult, s.Influence, rayResult.IsRecovered, m)
		} else if !rayResult.ApproachedBoundingBox {
			// Optimise the outside-bounding-box cases by skipping all further samples
			break
		}
	}
}

func setResult(result *RenderSample, element voxelobject.ProcessedElement, lighting geometry.Vector3, depth int, shadowLength int, influence float64, isRecovered bool, m manifest.Manifest) {

	if shadowLength > 0 && shadowLength < 10 {
		result.Shadowing = 1.0
	} else if shadowLength > 0 && shadowLength < 80 {
		result.Shadowing = float64(70-(shadowLength-10)) / 80.0
	}

	result.Collision = true
	result.Index = element.Index
	result.Depth = depth
	result.LightAmount = getLightingValue(element.AveragedNormal, lighting)
	if result.LightAmount > m.ShadowThreshold {
		if m.SoftShadow {
			result.Shadowing = result.Shadowing * (result.LightAmount - m.ShadowThreshold) / (1.0 - m.ShadowThreshold)
		}
	} else {
		result.Shadowing = 0.0
	}
	result.Normal = element.Normal
	result.Occlusion = element.Occlusion
	result.AveragedNormal = element.AveragedNormal
	result.Detail = element.Detail
	result.Influence = influence
	result.Count = 1
	result.IsRecovered = isRecovered
}

func getLightingValue(normal, lighting geometry.Vector3) float64 {
	return normal.Dot(lighting)
}
