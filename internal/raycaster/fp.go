package raycaster

import (
	"github.com/mattkimber/gorender/internal/geometry"
	"github.com/mattkimber/gorender/internal/voxelobject"
	"math"
	"math/rand"
)

var jitter []geometry.Vector3

func castFpRay(object voxelobject.ProcessedVoxelObject, loc0 geometry.Vector3, loc geometry.Vector3, ray geometry.Vector3, limits geometry.Vector3, flipY bool) (result RayResult) {
	if len(jitter) == 0 {
		jitter = make([]geometry.Vector3, 20)
		for i := 0; i < 20; i++ {
			jitter[i] = geometry.Vector3{X: rand.Float64(), Y: rand.Float64(), Z: rand.Float64()}.Normalise().MultiplyByConstant(0.01)
		}
	}

	if collision, loc := castRayToCandidate(object, loc, ray, limits, flipY); collision {
		lx, ly, lz := recoverNonSurfaceVoxel(object, loc, ray, limits, flipY)
		return RayResult{X: lx, Y: ly, Z: lz, HasGeometry: true, Depth: int(loc0.Subtract(loc).Length())}
	}

	return
}

func castRayToCandidate(object voxelobject.ProcessedVoxelObject, loc geometry.Vector3, ray geometry.Vector3, limits geometry.Vector3, flipY bool) (bool, geometry.Vector3) {
	i, fi := 0, 0.0
	bSizeY := uint8(object.Size.Y - 1)
	loc0 := loc

	for {
		// CanTerminate is an expensive check but we don't need to run it every cycle
		if i%4 == 0 && canTerminateRay(loc, ray, limits) {
			break
		}

		if isInsideBoundingVolume(loc, limits) {
			lx, ly, lz := byte(loc.X), byte(loc.Y), byte(loc.Z)

			if flipY {
				ly = bSizeY - ly
			}

			if object.Elements[lx][ly][lz].Index != 0 {
				return true, loc
			}
		}

		i++
		fi++
		loc = loc0.Add(ray.MultiplyByConstant(fi))
	}

	return false, geometry.Vector3{}
}

// Attempt to recover a non-surface voxel by taking a more DDA-like approach where we trace backward up the ray
// starting with X, then Y, then Z, then repeat until we find a surface voxel or bail.
func recoverNonSurfaceVoxel(object voxelobject.ProcessedVoxelObject, loc geometry.Vector3, ray geometry.Vector3, limits geometry.Vector3, flipY bool) (lx byte, ly byte, lz byte) {

	bSizeY := uint8(object.Size.Y - 1)

	lx, ly, lz = byte(loc.X), byte(loc.Y), byte(loc.Z)
	if flipY {
		ly = bSizeY - ly
	}

	if isInsideBoundingVolume(loc, limits) && object.Elements[lx][ly][lz].IsSurface {
		return
	}

	// Check always checks a 9 voxel "halo"
	check := make([]geometry.PointB, 9)
	checkOrder := []int{4, 1, 7, 3, 5, 0, 2, 6, 8}

	loc0 := loc
	x, y, z := ray.X, ray.Y, ray.Z

	for i := 0; i < 10; i++ {

		loc = loc.Subtract(ray.Normalise())
		lx, ly, lz = byte(loc.X), byte(loc.Y), byte(loc.Z)
		if flipY {
			ly = bSizeY - ly
		}

		for j := 0; j < 3; j++ {

			if math.Abs(x) > math.Abs(y) && math.Abs(x) > math.Abs(z) {
				// X-major

				for k := byte(0); k < 9; k++ {
					check[k] = geometry.PointB{X: lx, Y: ly - 1 + (k % 3), Z: lz - 1 + (k / 3)}
				}

				x = 0
			} else if math.Abs(y) > math.Abs(x) && math.Abs(y) > math.Abs(z) {
				// Y-major

				for k := byte(0); k < 9; k++ {
					check[k] = geometry.PointB{X: lx - 1 + (k % 3), Y: ly, Z: lz - 1 + (k / 3)}
				}

				y = 0
			} else if math.Abs(z) > math.Abs(x) && math.Abs(z) > math.Abs(y) {
				// Z-major

				for k := byte(0); k < 9; k++ {
					check[k] = geometry.PointB{X: lx - 1 + (k % 3), Y: ly - 1 + (k / 3), Z: lz}
				}

				z = 0
			}

			for k := byte(0); k < 9; k++ {
				point := check[checkOrder[k]]
				pointF := geometry.Vector3{X: float64(point.X), Y: float64(point.Y), Z: float64(point.Z)}

				lx, ly, lz = point.X, point.Y, point.Z
				if flipY {
					ly = bSizeY - ly
				}

				if isInsideBoundingVolume(pointF, limits) {
					if object.Elements[lx][ly][lz].IsSurface {
						//fmt.Printf("Recovered surface voxel at %d %d %d (%d)\n", lx, ly, lz, k)
						//fmt.Printf("%v\n", check)
						return
					}
				}
			}

			if x == 0 && y == 0 && z == 0 {
				x, y, z = ray.X, ray.Y, ray.Z
			}
		}

	}

	lx, ly, lz = byte(loc0.X), byte(loc0.Y), byte(loc0.Z)

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

func isInsideBoundingVolume(loc geometry.Vector3, limits geometry.Vector3) bool {
	return loc.X >= 0 && loc.Y >= 0 && loc.Z >= 0 && loc.X < limits.X && loc.Y < limits.Y && loc.Z < limits.Z
}

func canTerminateRay(loc geometry.Vector3, ray geometry.Vector3, limits geometry.Vector3) bool {
	return (loc.X < 0 && ray.X <= 0) || (loc.Y < 0 && ray.Y <= 0) || (loc.Z < 0 && ray.Z <= 0) ||
		(loc.X > limits.X && ray.X >= 0) || (loc.Y > limits.Y && ray.Y >= 0) || (loc.Z > limits.Z && ray.Z >= 0)
}
