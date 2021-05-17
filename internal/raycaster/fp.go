package raycaster

import (
	"github.com/mattkimber/gorender/internal/geometry"
	"github.com/mattkimber/gorender/internal/voxelobject"
	"math"
)

func castFpRay(object voxelobject.ProcessedVoxelObject, loc0 geometry.Vector3, loc geometry.Vector3, ray geometry.Vector3, limits geometry.Vector3, flipY bool) (result RayResult) {
	if collision, loc, approachedBB := castRayToCandidate(object, loc, ray, limits, flipY); collision {
		lx, ly, lz, isRecovered := recoverNonSurfaceVoxel(object, loc, ray, limits, flipY)
		return RayResult{
			X:                     lx,
			Y:                     ly,
			Z:                     lz,
			IsRecovered:           isRecovered,
			HasGeometry:           true,
			Depth:                 int(loc0.Subtract(loc).Length()),
			ApproachedBoundingBox: approachedBB,
		}
	} else if approachedBB {
		return RayResult{ApproachedBoundingBox: true}
	}

	return
}

func castRayToCandidate(object voxelobject.ProcessedVoxelObject, loc geometry.Vector3, ray geometry.Vector3, limits geometry.Vector3, flipY bool) (bool, geometry.Vector3, bool) {
	i, fi := 0, 0.0
	bSizeY := object.Size.Y - 1
	loc0 := loc
	approachedBB := false

	for {
		// CanTerminate is an expensive check but we don't need to run it every cycle
		if i%4 == 0 && canTerminateRay(loc, ray, limits) {
			break
		}

		if isInsideBoundingVolume(loc, limits) {
			approachedBB = true
			lx, ly, lz := int(loc.X), int(loc.Y), int(loc.Z)

			if flipY {
				ly = bSizeY - ly
			}

			if object.Elements[lx][ly][lz].Index != 0 {
				return true, loc, approachedBB
			}
		} else if !approachedBB && isNearlyInsideBoundingVolume(loc, limits) {
			approachedBB = true
		}

		i++
		fi++
		loc = loc0.Add(ray.MultiplyByConstant(fi))
	}

	return false, geometry.Vector3{}, approachedBB
}

// Attempt to recover a non-surface voxel by taking a more DDA-like approach where we trace backward up the ray
// starting with X, then Y, then Z, then repeat until we find a surface voxel or bail.
func recoverNonSurfaceVoxel(object voxelobject.ProcessedVoxelObject, loc geometry.Vector3, ray geometry.Vector3, limits geometry.Vector3, flipY bool) (lx int, ly int, lz int, isRecovered bool) {

	bSizeY := object.Size.Y - 1

	lx, ly, lz = int(loc.X), int(loc.Y), int(loc.Z)
	if flipY {
		ly = bSizeY - ly
	}

	if isInsideBoundingVolume(loc, limits) && object.Elements[lx][ly][lz].IsSurface {
		return
	}

	// Signify this voxel was recovered
	isRecovered = true

	// Check always checks a 9 voxel "halo"
	check := make([]geometry.Point, 9)
	checkOrder := []int{4, 1, 7, 3, 5, 0, 2, 6, 8}

	loc0 := loc
	x, y, z := ray.X, ray.Y, ray.Z

	for i := 0; i < 10; i++ {
		lx, ly, lz = int(loc.X), int(loc.Y), int(loc.Z)
		if flipY {
			ly = bSizeY - ly
		}

		for j := 0; j < 3; j++ {

			if math.Abs(x) > math.Abs(y) && math.Abs(x) > math.Abs(z) {
				// X-major

				for k := 0; k < 9; k++ {
					check[k] = geometry.Point{X: lx, Y: ly - 1 + (k % 3), Z: lz - 1 + (k / 3)}
				}

				x = 0
			} else if math.Abs(y) > math.Abs(x) && math.Abs(y) > math.Abs(z) {
				// Y-major

				for k := 0; k < 9; k++ {
					check[k] = geometry.Point{X: lx - 1 + (k % 3), Y: ly, Z: lz - 1 + (k / 3)}
				}

				y = 0
			} else if math.Abs(z) > math.Abs(x) && math.Abs(z) > math.Abs(y) {
				// Z-major

				for k := 0; k < 9; k++ {
					check[k] = geometry.Point{X: lx - 1 + (k % 3), Y: ly - 1 + (k / 3), Z: lz}
				}

				z = 0
			}

			for k := 0; k < 9; k++ {
				point := check[checkOrder[k]]
				pointF := geometry.Vector3{X: float64(point.X), Y: float64(point.Y), Z: float64(point.Z)}

				lx, ly, lz = point.X, point.Y, point.Z

				if isInsideBoundingVolume(pointF, limits) {
					if object.Elements[lx][ly][lz].IsSurface {
						return
					}
				}
			}

			if x == 0 && y == 0 && z == 0 {
				x, y, z = ray.X, ray.Y, ray.Z
			}
		}

		loc = loc.Subtract(ray.Normalise())
	}

	lx, ly, lz = int(loc0.X), int(loc0.Y), int(loc0.Z)

	if flipY {
		ly = bSizeY - ly
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

func isNearlyInsideBoundingVolume(loc geometry.Vector3, limits geometry.Vector3) bool {
	// We are within 3 voxels of the bounding box, which is considered "approached"
	return loc.X >= -3 && loc.Y >= -3 && loc.Z >= -3 && loc.X < limits.X + 3 && loc.Y < limits.Y + 3 && loc.Z < limits.Z + 3
}

func isInsideBoundingVolume(loc geometry.Vector3, limits geometry.Vector3) bool {
	return loc.X >= 0 && loc.Y >= 0 && loc.Z >= 0 && loc.X < limits.X && loc.Y < limits.Y && loc.Z < limits.Z
}

func canTerminateRay(loc geometry.Vector3, ray geometry.Vector3, limits geometry.Vector3) bool {
	return (loc.X < 0 && ray.X <= 0) || (loc.Y < 0 && ray.Y <= 0) || (loc.Z < 0 && ray.Z <= 0) ||
		(loc.X > limits.X && ray.X >= 0) || (loc.Y > limits.Y && ray.Y >= 0) || (loc.Z > limits.Z && ray.Z >= 0)
}
