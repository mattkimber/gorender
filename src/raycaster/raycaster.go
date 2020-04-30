package raycaster

import (
	"geometry"
	"math"
	"voxelobject"
)

type RenderInfo struct {
	Collision              bool
	Index                  byte
	Normal, AveragedNormal geometry.Vector3
	Depth                  int
	LightAmount            float64
}

const lightingAngle = 60

type RenderOutput [][]RenderInfo

func getRenderDirection(angle int) geometry.Vector3 {
	x, y, z := -math.Cos(degToRad(angle)), math.Sin(degToRad(angle)), math.Sin(degToRad(30))
	return geometry.Vector3{X: x, Y: y, Z: z}.Normalise()
}

func degToRad(angle int) float64 {
	return (float64(angle) / 180.0) * math.Pi
}

func getLightingDirection(angle int) geometry.Vector3 {
	x, y, z := -math.Cos(degToRad(angle)), math.Sin(degToRad(angle)), math.Sin(degToRad(-45))
	return geometry.Zero().Subtract(geometry.Vector3{X: x, Y: y, Z: z}).Normalise()
}

func getViewportPlane(angle int, size geometry.Point) geometry.Plane {
	midpoint := geometry.Vector3{X: float64(size.X) / 2.0, Y: float64(size.Y) / 2.0, Z: float64(size.Z) - (float64(size.Y) / 2.0)}
	viewpoint := midpoint.Add(getRenderDirection(angle).MultiplyByConstant(100.0))

	planeNormal := geometry.UnitZ().MultiplyByConstant(midpoint.X)
	renderNormal := getRenderNormal(angle).MultiplyByConstant(midpoint.X)

	a := viewpoint.Subtract(renderNormal).Add(planeNormal)
	b := viewpoint.Add(renderNormal).Add(planeNormal)
	c := viewpoint.Add(renderNormal).Subtract(planeNormal)
	d := viewpoint.Subtract(renderNormal).Subtract(planeNormal)

	return geometry.Plane{A: a, B: b, C: c, D: d}
}

func getRenderNormal(angle int) geometry.Vector3 {
	x, y := -math.Cos(degToRad(angle)), math.Sin(degToRad(angle))
	return geometry.Vector3{X: y, Y: -x}.Normalise()
}

func isInsideBoundingVolume(loc geometry.Vector3, limits geometry.Vector3) bool {
	return loc.X >= 0 && loc.Y >= 0 && loc.Z >= 0 && loc.X < limits.X && loc.Y < limits.Y && loc.Z < limits.Z
}

func canTerminateRay(loc geometry.Vector3, ray geometry.Vector3, limits geometry.Vector3) bool {
	return (loc.X < 0 && ray.X <= 0) || (loc.Y < 0 && ray.Y <= 0) || (loc.Z < 0 && ray.Z <= 0) ||
		(loc.X > limits.X && ray.X >= 0) || (loc.Y > limits.Y && ray.Y >= 0) || (loc.Z > limits.Z && ray.Z >= 0)
}

func GetRaycastOutput(object voxelobject.ProcessedVoxelObject, angle int, w int, h int, debug bool) RenderOutput {
	size := object.Size

	limits := geometry.Vector3{X: float64(size.X), Y: float64(size.Y), Z: float64(size.Z)}

	viewport := getViewportPlane(angle, size)
	ray := geometry.Zero().Subtract(getRenderDirection(angle)).MultiplyByConstant(0.5)

	lighting := getLightingDirection(angle + lightingAngle)
	result := make(RenderOutput, w)

	for x := 0; x < w; x++ {
		result[x] = make([]RenderInfo, h)
		for y := 0; y < h; y++ {
			u, v := float64(x)/float64(w), float64(y)/float64(h)
			loc := viewport.BiLerpWithinPlane(u, v)
			depth := 0

			for {
				if canTerminateRay(loc, ray, limits) {
					break
				}

				if isInsideBoundingVolume(loc, limits) {
					lx, ly, lz := byte(loc.X), byte(loc.Y), byte(loc.Z)
					if object.Elements[lx][ly][lz].Index != 0 {
						result[x][y].Collision = true
						result[x][y].Index = object.Elements[lx][ly][lz].Index
						result[x][y].Depth = depth
						result[x][y].LightAmount = getLightingValue(object.Elements[lx][ly][lz].AveragedNormal, lighting)
						if debug {
							result[x][y].Normal = object.Elements[lx][ly][lz].Normal
							result[x][y].AveragedNormal = object.Elements[lx][ly][lz].AveragedNormal
						}
					}
				}

				loc = loc.Add(ray)
				depth++
			}
		}
	}

	return result
}

func getLightingValue(normal, lighting geometry.Vector3) float64 {
	return normal.Dot(lighting)
}
