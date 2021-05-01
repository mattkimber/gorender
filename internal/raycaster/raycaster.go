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
	Influence 			   float64
	Detail				   float64
	IsRecovered            bool
}

type RayResult struct {
	X, Y, Z     int
	HasGeometry bool
	Depth       int
	IsRecovered bool
	HitBoundingBox bool
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
		if maxX > 255 {
			maxX = 255
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
				raycastSamples(viewport, &samples, ray, limits, object, spr, lighting, result, thisX, y, minX, maxX, joggle)
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
	spr manifest.Sprite,
	lighting geometry.Vector3,
	result RenderOutput,
	thisX int,
	y int,
	minX int,
	maxX int,
	joggle float64) {
	for i, s := range *samples {
		loc0 := viewport.BiLerpWithinPlane(s.Location.X, s.Location.Y)
		loc0.Z += joggle
		loc := getIntersectionWithBounds(loc0, ray, limits)

		rayResult := castFpRay(object, loc0, loc, ray, limits, spr.Flip)
		if rayResult.HasGeometry && rayResult.X >= minX && rayResult.X <= maxX {
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
			shadowResult := castFpRay(object, shadowLoc, shadowLoc, shadowVec, limits, false).Depth
			setResult(&result[thisX][y][i], object.Elements[rayResult.X][rayResult.Y][rayResult.Z], lighting, rayResult.Depth, shadowResult, s.Influence, rayResult.IsRecovered)
		} else if !rayResult.HitBoundingBox {
			// Optimise the outside-bounding-box cases by skipping all further samples
			break
		}
	}

}

func setResult(result *RenderSample, element voxelobject.ProcessedElement, lighting geometry.Vector3, depth int, shadowLength int, influence float64, isRecovered bool) {

	if shadowLength > 0 && shadowLength < 10 {
		result.Shadowing = 1.0
	} else if shadowLength > 0 && shadowLength < 80 {
		result.Shadowing = float64(70-(shadowLength-10)) / 80.0
	}

	result.Collision = true
	result.Index = element.Index
	result.Depth = depth
	result.LightAmount = getLightingValue(element.AveragedNormal, lighting)
	result.Normal = element.Normal
	result.Occlusion = element.Occlusion
	result.AveragedNormal = element.AveragedNormal
	result.Detail = element.Detail
	result.Influence = influence
	result.IsRecovered = isRecovered
}

func getLightingValue(normal, lighting geometry.Vector3) float64 {
	return normal.Dot(lighting)
}
