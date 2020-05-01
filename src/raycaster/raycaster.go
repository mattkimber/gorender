package raycaster

import (
	"geometry"
	"voxelobject"
)

type RenderInfo struct {
	Collision              bool
	Index                  byte
	Normal, AveragedNormal geometry.Vector3
	Depth                  int
	LightAmount            float64
}

type RayResult struct {
	X,Y,Z byte
	HasGeometry bool
	Depth int
}

const lightingAngle = 60

type RenderOutput [][]RenderInfo

func GetRaycastOutput(object voxelobject.ProcessedVoxelObject, angle int, w int, h int) RenderOutput {
	size := object.Size

	limits := geometry.Vector3{X: float64(size.X), Y: float64(size.Y), Z: float64(size.Z)}

	viewport := getViewportPlane(angle, size)
	ray := geometry.Zero().Subtract(getRenderDirection(angle))

	lighting := getLightingDirection(angle + lightingAngle)
	result := make(RenderOutput, w)

	for x := 0; x < w; x++ {
		result[x] = make([]RenderInfo, h)
		for y := 0; y < h; y++ {
			rayResult := castFpRay(object, float64(x)/float64(w), float64(y)/float64(h), viewport, ray, limits)
			if rayResult.HasGeometry {
				setResult(&result[x][y], object.Elements[rayResult.X][rayResult.Y][rayResult.Z], lighting, rayResult.Depth)
			}
		}
	}

	return result
}

func setResult(result *RenderInfo, element voxelobject.ProcessedElement, lighting geometry.Vector3, depth int) {
	result.Collision = true
	result.Index = element.Index
	result.Depth = depth
	result.LightAmount = getLightingValue(element.AveragedNormal, lighting)
	result.Normal = element.Normal
	result.AveragedNormal = element.AveragedNormal
}

func getLightingValue(normal, lighting geometry.Vector3) float64 {
	return normal.Dot(lighting)
}
