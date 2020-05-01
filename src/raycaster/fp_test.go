package raycaster

import (
	"geometry"
	"testing"
)

func TestIsInsideBoundingVolume(t *testing.T) {
	testCases := []struct {
		loc, limits geometry.Vector3
		expected    bool
	}{
		{geometry.Vector3{}, geometry.Vector3{X: 2, Y: 2, Z: 2}, true},
		{geometry.Vector3{X: 1, Y: 1, Z: 1}, geometry.Vector3{X: 2, Y: 2, Z: 2}, true},
		{geometry.Vector3{X: 4, Y: 2, Z: 2}, geometry.Vector3{X: 2, Y: 2, Z: 2}, false},
		{geometry.Vector3{X: -1, Y: 2, Z: 2}, geometry.Vector3{X: 2, Y: 2, Z: 2}, false},
	}

	for _, testCase := range testCases {
		if result := isInsideBoundingVolume(testCase.loc, testCase.limits); result != testCase.expected {
			t.Errorf("co-ordinates %v inside %v expected %v, got %v", testCase.loc, testCase.limits, testCase.expected, result)
		}
	}
}

func TestCanTerminateRay(t *testing.T) {
	testCases := []struct {
		loc, ray, limits geometry.Vector3
		expected         bool
	}{
		{geometry.Vector3{X: 1}, geometry.Vector3{X: 1}, geometry.Vector3{X: 2, Y: 2, Z: 2}, false},
		{geometry.Vector3{X: 1}, geometry.Vector3{X: -1}, geometry.Vector3{X: 2, Y: 2, Z: 2}, false},
		{geometry.Vector3{X: -1}, geometry.Vector3{X: 1}, geometry.Vector3{X: 2, Y: 2, Z: 2}, false},
		{geometry.Vector3{X: -1}, geometry.Vector3{X: -1}, geometry.Vector3{X: 2, Y: 2, Z: 2}, true},
		{geometry.Vector3{X: 3}, geometry.Vector3{X: 1}, geometry.Vector3{X: 2, Y: 2, Z: 2}, true},
		{geometry.Vector3{X: 3}, geometry.Vector3{X: -1}, geometry.Vector3{X: 2, Y: 2, Z: 2}, false},
	}

	for _, testCase := range testCases {
		if result := canTerminateRay(testCase.loc, testCase.ray, testCase.limits); result != testCase.expected {
			t.Errorf("co-ordinates %v can terminate ray %v for limits %v expected %v, got %v", testCase.loc, testCase.ray, testCase.limits, testCase.expected, result)
		}
	}
}

func Test_castFpRay(t *testing.T) {
	object := getObject("testcube", t)
	size := object.Size
	limits := geometry.Vector3{X: float64(size.X), Y: float64(size.Y), Z: float64(size.Z)}

	viewport := geometry.Plane{
		A: geometry.Vector3{X: 8, Z: 5},
		B: geometry.Vector3{X: 8, Y: 4, Z: 5},
		C: geometry.Vector3{X: 8, Y: 4, Z: 1},
		D: geometry.Vector3{X: 8, Z: 1},
	}

	ray := geometry.Vector3{X: -1, Y: 0, Z: -0.125}.Normalise()
	result := castFpRay(object, 0.5, 0.5, viewport, ray, limits)

	if !result.HasGeometry {
		t.Errorf("did not find geometry")
	}

	if result.X != 2 || result.Y != 2 || result.Z != 2 {
		t.Errorf("incorrect voxel - expected [2,2,2], got [%d,%d,%d]", result.X, result.Y, result.Z)
	}

	if result.Depth != 6 {
		t.Errorf("incorrect depth - expected 6, got %d", result.Depth)
	}
}
