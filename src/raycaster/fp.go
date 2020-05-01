package raycaster

import (
	"geometry"
	"voxelobject"
)

func castFpRay(object voxelobject.ProcessedVoxelObject, u, v float64, viewport geometry.Plane, ray geometry.Vector3, limits geometry.Vector3) (result RayResult) {
	loc := viewport.BiLerpWithinPlane(u, v)
	depth := 0

	for {
		if canTerminateRay(loc, ray, limits) {
			break
		}

		if isInsideBoundingVolume(loc, limits) {
			lx, ly, lz := byte(loc.X), byte(loc.Y), byte(loc.Z)
			if object.Elements[lx][ly][lz].Index != 0 {
				return RayResult{X:lx, Y:ly, Z:lz, HasGeometry: true, Depth: depth}
			}
		}

		loc = loc.Add(ray)
		depth++
	}

	return
}

func isInsideBoundingVolume(loc geometry.Vector3, limits geometry.Vector3) bool {
	return loc.X >= 0 && loc.Y >= 0 && loc.Z >= 0 && loc.X < limits.X && loc.Y < limits.Y && loc.Z < limits.Z
}

func canTerminateRay(loc geometry.Vector3, ray geometry.Vector3, limits geometry.Vector3) bool {
	return (loc.X < 0 && ray.X <= 0) || (loc.Y < 0 && ray.Y <= 0) || (loc.Z < 0 && ray.Z <= 0) ||
		(loc.X > limits.X && ray.X >= 0) || (loc.Y > limits.Y && ray.Y >= 0) || (loc.Z > limits.Z && ray.Z >= 0)
}
