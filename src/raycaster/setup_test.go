package raycaster

import (
	"geometry"
	"testing"
)

func TestGetRenderDirection(t *testing.T) {
	testCases := []struct {
		angle    int
		expected geometry.Vector3
	}{
		{0, geometry.Vector3{X: -0.894427190999916, Z: 0.447213595499958}},
		{45, geometry.Vector3{X: -0.632455532033676, Y: 0.632455532033676, Z: 0.447213595499958}},
		{90, geometry.Vector3{Y: 0.894427190999916, Z: 0.447213595499958}},
	}

	for _, testCase := range testCases {
		if result := getRenderDirection(testCase.angle); !result.Equals(testCase.expected) {
			t.Errorf("Angle %d expected render direction %v, got %v", testCase.angle, testCase.expected, result)
		}
	}
}

func TestGetViewportPlane(t *testing.T) {
	testCases := []struct {
		angle    int
		x, y     int
		expected geometry.Plane
	}{
		{0, 126, 40, geometry.Plane{
			A: geometry.Vector3{X: -26.442719099991592, Y: -43, Z: 1.721359549995782},
			B: geometry.Vector3{X: -26.442719099991592, Y: 83, Z: 1.721359549995782},
			C: geometry.Vector3{X: -26.442719099991592, Y: 83, Z: 127.72135954999578},
			D: geometry.Vector3{X: -26.442719099991592, Y: -43, Z: 127.72135954999578},
		}},
		{45, 126, 40, geometry.Plane{
			A: geometry.Vector3{X: -44.79328041812008, Y: 38.69782598861508, Z: 1.721359549995782},
			B: geometry.Vector3{X: 44.302174011384906, Y: 127.79328041812008, Z: 1.721359549995782},
			C: geometry.Vector3{X: 44.302174011384906, Y: 127.79328041812008, Z: 127.72135954999578},
			D: geometry.Vector3{X: -44.79328041812008, Y: 38.69782598861508, Z: 127.72135954999578},
		}},
		{90, 126, 40, geometry.Plane{
			A: geometry.Vector3{Y: 109.44271909999159, Z: 1.721359549995782},
			B: geometry.Vector3{X: 126, Y: 109.44271909999159, Z: 1.721359549995782},
			C: geometry.Vector3{X: 126, Y: 109.44271909999159, Z: 127.72135954999578},
			D: geometry.Vector3{Y: 109.44271909999159, Z: 127.72135954999578},
		}},
	}

	for _, testCase := range testCases {
		if result := getViewportPlane(testCase.angle, geometry.Point{testCase.x, testCase.y, testCase.y}); !result.Equals(testCase.expected) {
			t.Errorf("Angle %d expected viewport plane %v, got %v", testCase.angle, testCase.expected, result)
		}
	}
}
