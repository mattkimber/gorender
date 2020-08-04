package raycaster

import (
	"github.com/mattkimber/gorender/internal/geometry"
	"github.com/mattkimber/gorender/internal/manifest"
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
			A: geometry.Vector3{X: -63, Y: 0, Z: 0},
			B: geometry.Vector3{X: -63, Y: 40, Z: 0},
			C: geometry.Vector3{X: -63, Y: 40, Z: 40},
			D: geometry.Vector3{X: -63, Y: 0, Z: 40},
		}},
		{90, 126, 40, geometry.Plane{
			A: geometry.Vector3{X: 0, Y: 146, Z: 0},
			B: geometry.Vector3{X: 126, Y: 146, Z: 0},
			C: geometry.Vector3{X: 126, Y: 146, Z: 40},
			D: geometry.Vector3{X: 0, Y: 146, Z: 40},
		}},
		{180, 126, 40, geometry.Plane{
			A: geometry.Vector3{X: 189, Y: 40, Z: 0},
			B: geometry.Vector3{X: 189, Y: 0, Z: 0},
			C: geometry.Vector3{X: 189, Y: 0, Z: 40},
			D: geometry.Vector3{X: 189, Y: 40, Z: 40},
		}},
		{270, 126, 40, geometry.Plane{
			A: geometry.Vector3{X: 126, Y: -106, Z: 0},
			B: geometry.Vector3{X: 0, Y: -106, Z: 0},
			C: geometry.Vector3{X: 0, Y: -106, Z: 40},
			D: geometry.Vector3{X: 126, Y: -106, Z: 40},
		}},
	}

	for _, testCase := range testCases {
		size := geometry.Point{X: testCase.x, Y: testCase.y, Z: testCase.y}
		mSize := geometry.Vector3{X: float64(testCase.x), Y: float64(testCase.y), Z: float64(testCase.y)}
		m := manifest.Manifest{Size: mSize}
		if result := getViewportPlane(testCase.angle, m, 0, size); !result.Equals(testCase.expected) {
			t.Errorf("Angle %f expected viewport plane %v, got %v", testCase.angle, testCase.expected, result)
		}
	}
}
