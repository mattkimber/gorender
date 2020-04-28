package raycaster

import (
	"geometry"
	"math"
	"voxelobject"
)

type RenderInfo struct {
	Collision bool
	Index     byte
}

type RenderOutput [][]RenderInfo

func getRenderDirection(angle int) geometry.Vector3 {
	x, y, z := -math.Cos(degToRad(angle)), math.Sin(degToRad(angle)), math.Sin(degToRad(30))
	return geometry.Vector3{X: x, Y: y, Z: z}.Normalise()
}

func degToRad(angle int) float64 {
	return (float64(angle) / 180.0) * math.Pi
}

func getViewportPlane(angle int, x byte, y byte) geometry.Plane {
	midpoint := geometry.Vector3{X: float64(x) / 2.0, Y: float64(y) / 2.0, Z: float64(y) / 2.0}
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

func GetRaycastOutput(object voxelobject.RawVoxelObject, angle int, w int, h int) RenderOutput {
	size := object.Size()

	limits := geometry.Vector3{X: float64(size.X), Y: float64(size.Y), Z: float64(size.Z)}

	viewport := getViewportPlane(angle, size.X, size.Y)
	ray := geometry.Zero().Subtract(getRenderDirection(angle))

	result := make(RenderOutput, w)

	for x := 0; x < w; x++ {
		result[x] = make([]RenderInfo, h)
		for y := 0; y < h; y++ {
			u, v := float64(x)/float64(w), float64(y)/float64(h)
			loc := viewport.BiLerpWithinPlane(u, v)

			for {
				if canTerminateRay(loc, ray, limits) {
					break
				}

				if isInsideBoundingVolume(loc, limits) {
					lx, ly, lz := byte(loc.X), byte(loc.Y), byte(loc.Z)
					if object[lx][ly][lz] != 0 {
						result[x][y].Collision = true
						result[x][y].Index = object[lx][ly][lz]
					}
				}

				loc = loc.Add(ray)
			}
		}
	}

	return result
}
