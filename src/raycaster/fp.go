package raycaster

import (
	"geometry"
	"math/rand"
	"voxelobject"
)

func castFpRay(object voxelobject.ProcessedVoxelObject, u, v float64, viewport geometry.Plane, ray geometry.Vector3, limits geometry.Vector3) (result RayResult) {
	loc0 := viewport.BiLerpWithinPlane(u, v)
	loc := getIntersectionWithBounds(loc0, ray, limits)

	for {
		if canTerminateRay(loc, ray, limits) {
			break
		}

		if isInsideBoundingVolume(loc, limits) {
			lx, ly, lz := byte(loc.X), byte(loc.Y), byte(loc.Z)
			if object.Elements[lx][ly][lz].Index != 0 {
				lx, ly, lz = recoverNonSurfaceVoxel(object, loc, ray, limits)
				return RayResult{X: lx, Y: ly, Z: lz, HasGeometry: true, Depth: int(loc0.Subtract(loc).Length())}
			}
		}

		loc = loc.Add(ray)
	}

	return
}

func recoverNonSurfaceVoxel(object voxelobject.ProcessedVoxelObject, loc geometry.Vector3, ray geometry.Vector3, limits geometry.Vector3) (lx byte, ly byte, lz byte) {
	if !object.Elements[lx][ly][lz].IsSurface {
		loc2 := loc
		lx, ly, lz = byte(loc.X), byte(loc.Y), byte(loc.Z)

		for i := 0; i < 20; i++ {
			loc2 = loc2.Subtract(ray.MultiplyByConstant(0.125))
			if !isInsideBoundingVolume(loc2, limits) {
				break
			}

			lx, ly, lz = byte(loc2.X), byte(loc2.Y), byte(loc2.Z)
			if object.Elements[lx][ly][lz].IsSurface {
				break
			}
		}
	}

	if !object.Elements[lx][ly][lz].IsSurface {
		loc2 := loc

		for i := 0; i < 20; i++ {
			jitter := geometry.Vector3{X: rand.Float64(), Y: rand.Float64(), Z: rand.Float64()}.Normalise().MultiplyByConstant(0.01)
			loc2 = loc2.Subtract(ray.MultiplyByConstant(0.125).Add(jitter))
			if !isInsideBoundingVolume(loc2, limits) {
				break
			}

			lx, ly, lz = byte(loc2.X), byte(loc2.Y), byte(loc2.Z)
			if object.Elements[lx][ly][lz].IsSurface {
				break
			}
		}

		lx, ly, lz = byte(loc.X), byte(loc.Y), byte(loc.Z)
	}
	return
}

func getIntersectionWithBounds(loc, ray, limits geometry.Vector3) geometry.Vector3 {
	if canTerminateRay(loc, ray, limits) {
		return loc
	}

	loc = loc.Add(getIntersectionVector(ray.X, loc.X, limits.X, ray))
	loc = loc.Add(getIntersectionVector(ray.Y, loc.Y, limits.Y, ray))

	return loc
}

func getIntersectionVector(rayDimension, locDimension, limitDimension float64, ray geometry.Vector3) geometry.Vector3 {
	dist := -1.0

	if rayDimension > 0.1 {
		dist = -locDimension
	}
	if rayDimension < -0.1 {
		dist = limitDimension - locDimension
	}

	if dist > 0 {
		return ray.MultiplyByConstant(dist / rayDimension)
	}

	return geometry.Zero()
}

func isInsideBoundingVolume(loc geometry.Vector3, limits geometry.Vector3) bool {
	return loc.X >= 0 && loc.Y >= 0 && loc.Z >= 0 && loc.X < limits.X && loc.Y < limits.Y && loc.Z < limits.Z
}

func canTerminateRay(loc geometry.Vector3, ray geometry.Vector3, limits geometry.Vector3) bool {
	return (loc.X < 0 && ray.X <= 0) || (loc.Y < 0 && ray.Y <= 0) || (loc.Z < 0 && ray.Z <= 0) ||
		(loc.X > limits.X && ray.X >= 0) || (loc.Y > limits.Y && ray.Y >= 0) || (loc.Z > limits.Z && ray.Z >= 0)
}
