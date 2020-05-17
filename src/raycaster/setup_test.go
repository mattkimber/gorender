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
			A: geometry.Vector3{X: -37, Y: 0, Z: -43},
			B: geometry.Vector3{X: -37, Y: 40, Z: -43},
			C: geometry.Vector3{X: -37, Y: 40, Z: 83},
			D: geometry.Vector3{X: -37, Y: 0, Z: 83},
		}},
		{45, 126, 40, geometry.Plane{
			A: geometry.Vector3{X: -49.21067811865475, Y: 49.21067811865474, Z: -38.68986283848345},
			B: geometry.Vector3{X: 33.78932188134524, Y: 132.21067811865476, Z: -38.68986283848345},
			C: geometry.Vector3{X: 33.78932188134524, Y: 132.21067811865476, Z: 78.68986283848345},
			D: geometry.Vector3{X: -49.21067811865475, Y: 49.21067811865474, Z: 78.68986283848345},
		}},
	}

	for _, testCase := range testCases {
		size := geometry.Point{X: testCase.x, Y: testCase.y, Z: testCase.y}
		m := manifest.Manifest{Size: size}
		if result := getViewportPlane(testCase.angle, m, size); !result.Equals(testCase.expected) {
			t.Errorf("Angle %f expected viewport plane %v, got %v", testCase.angle, testCase.expected, result)
		}
	}
}
