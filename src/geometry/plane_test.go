package geometry

import "testing"

func TestPlane_Equals(t *testing.T) {
	plane1 := Plane{A: Vector3{Y: 1}, B: Vector3{X: 1, Y: 1}, C: Vector3{X: 1}, D: Vector3{}}
	plane2 := Plane{A: Vector3{Y: 2}, B: Vector3{X: 1, Y: 2}, C: Vector3{X: 1}, D: Vector3{}}

	testCases := []struct {
		a, b     Plane
		expected bool
	}{
		{plane1, plane1, true},
		{plane2, plane2, true},
		{plane1, plane2, false},
	}

	for _, testCase := range testCases {
		if result := testCase.a.Equals(testCase.b); result != testCase.expected {
			t.Errorf("%v == %v expected %v, was %v", testCase.a, testCase.b, testCase.expected, result)
		}
	}
}

func TestPlane_BiLerpWithinPlane(t *testing.T) {
	plane := Plane{A: Vector3{Y: 2}, B: Vector3{X: 2, Y: 2}, C: Vector3{X: 2}, D: Vector3{}}

	testCases := []struct {
		p        Plane
		u, v     float64
		expected Vector3
	}{
		{plane, 0, 0, Vector3{Y: 2}},
		{plane, 1, 0, Vector3{X: 2, Y: 2}},
		{plane, 0, 1, Vector3{}},
		{plane, 1, 1, Vector3{X: 2}},
		{plane, 0.5, 0.5, Vector3{X: 1, Y: 1}},
	}

	for _, testCase := range testCases {
		if result := testCase.p.BiLerpWithinPlane(testCase.u, testCase.v); result != testCase.expected {
			t.Errorf("Bilerp %v with [%f,%f] expected %v, was %v", testCase.p, testCase.u, testCase.v, testCase.expected, result)
		}
	}
}
