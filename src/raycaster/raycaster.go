package raycaster

import (
	"geometry"
	"manifest"
	"sync"
	"voxelobject"
)

type RenderInfo struct {
	Collision              bool
	Index                  byte
	Normal, AveragedNormal geometry.Vector3
	Depth, Occlusion       int
	LightAmount            float64
	Shadowing              float64
}

type RayResult struct {
	X, Y, Z     byte
	HasGeometry bool
	Depth       int
}

type RenderOutput [][]RenderInfo

func GetRaycastOutput(object voxelobject.ProcessedVoxelObject, m manifest.Manifest, spr manifest.Sprite, w int, h int) RenderOutput {
	size := object.Size
	limits := geometry.Vector3{X: float64(size.X), Y: float64(size.Y), Z: float64(size.Z)}

	viewport := getViewportPlane(spr.Angle, m, size)
	ray := geometry.Zero().Subtract(getRenderDirection(spr.Angle, getElevationAngle(m)))

	lighting := getLightingDirection(spr.Angle+float64(m.LightingAngle), float64(m.LightingElevation), spr.Flip)
	result := make(RenderOutput, w)

	wg := sync.WaitGroup{}
	wg.Add(w)

	for x := 0; x < w; x++ {
		thisX := x
		go func() {
			result[thisX] = make([]RenderInfo, h)
			for y := 0; y < h; y++ {
				loc0 := viewport.BiLerpWithinPlane(float64(thisX)/float64(w), float64(y)/float64(h))
				loc := getIntersectionWithBounds(loc0, ray, limits)

				rayResult := castFpRay(object, loc0, loc, ray, limits, spr.Flip)
				if rayResult.HasGeometry {
					resultVec := geometry.Vector3{X: float64(rayResult.X), Y: float64(rayResult.Y), Z: float64(rayResult.Z)}
					shadowLoc := resultVec
					shadowVec := geometry.Zero().Subtract(lighting).Normalise()

					for {
						if !shadowLoc.Equals(resultVec) {
							break
						}

						shadowLoc = shadowLoc.Add(shadowVec)
					}

					shadowResult := castFpRay(object, shadowLoc, shadowLoc, shadowVec, limits, spr.Flip).Depth
					setResult(&result[thisX][y], object.Elements[rayResult.X][rayResult.Y][rayResult.Z], lighting, rayResult.Depth, shadowResult)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return result
}

func setResult(result *RenderInfo, element voxelobject.ProcessedElement, lighting geometry.Vector3, depth int, shadowLength int) {

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
}

func getLightingValue(normal, lighting geometry.Vector3) float64 {
	return normal.Dot(lighting)
}
