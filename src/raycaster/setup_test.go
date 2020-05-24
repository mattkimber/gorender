package raycaster

import (
	"geometry"
	"manifest"
	"testing"
)

func TestGetRenderDirection(t *testing.T) {
	testCases := []struct {
		angle    float64
		expected geometry.Vector3
	}{
		{0, geometry.Vector3{X: -0.894427190999916, Z: 0.447213595499958}},
		{45, geometry.Vector3{X: -0.632455532033676, Y: 0.632455532033676, Z: 0.447213595499958}},
		{90, geometry.Vector3{Y: 0.894427190999916, Z: 0.447213595499958}},
	}

	for _, testCase := range testCases {
		if result := getRenderDirection(testCase.angle, 30); !result.Equals(testCase.expected) {
			t.Errorf("Angle %f expected render direction %v, got %v", testCase.angle, testCase.expected, result)
		}
	}
}

func TestGetViewportPlane(t *testing.T) {
	testCases := []struct {
		angle    float64
		x, y     int
		expected geometry.Plane
	}{
		{0, 126, 40, geometry.Plane{
			A: geometry.Vector3{X: -37.5, Y: -0.5, Z: -0.5},
			B: geometry.Vector3{X: -37.5, Y: 39.5, Z: -0.5},
			C: geometry.Vector3{X: -37.5, Y: 39.5, Z: 39.5},
			D: geometry.Vector3{X: -37.5, Y: -0.5, Z: 39.5},
		}},
		{45, 126, 40, geometry.Plane{
			A: geometry.Vector3{X: -49.71067811865475, Y: 48.71067811865474, Z: -0.5},
			B: geometry.Vector3{X: 33.28932188134524, Y: 131.71067811865476, Z: -0.5},
			C: geometry.Vector3{X: 33.28932188134524, Y: 131.71067811865476, Z: 39.5},
			D: geometry.Vector3{X: -49.71067811865475, Y: 48.71067811865474, Z: 39.5},
		}},
	}

	for _, testCase := range testCases {
		size := geometry.Point{X: testCase.x, Y: testCase.y, Z: testCase.y}
		mSize := geometry.Vector3{X: float64(testCase.x), Y: float64(testCase.y), Z: float64(testCase.y)}
		m := manifest.Manifest{Size: mSize}
		if result := getViewportPlane(testCase.angle, m, size); !result.Equals(testCase.expected) {
			t.Errorf("Angle %f expected viewport plane %v, got %v", testCase.angle, testCase.expected, result)
		}
	}
}
