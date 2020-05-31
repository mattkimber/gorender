package geometry

import (
	"math"
	"testing"
)

func testVector(expected Vector3, vector Vector3, t *testing.T) {
	if vector != expected {
		t.Errorf("Expected %v, got %v", expected, vector)
	}
}

func testVector2(expected Vector2, vector Vector2, t *testing.T) {
	if vector != expected {
		t.Errorf("Expected %v, got %v", expected, vector)
	}
}

func TestZero(t *testing.T) {
	testVector(Vector3{0, 0, 0}, Zero(), t)
}

func TestUnitX(t *testing.T) {
	testVector(Vector3{1, 0, 0}, UnitX(), t)
}

func TestUnitY(t *testing.T) {
	testVector(Vector3{0, 1, 0}, UnitY(), t)
}

func TestUnitZ(t *testing.T) {
	testVector(Vector3{0, 0, 1}, UnitZ(), t)
}

func TestVector3_Add(t *testing.T) {
	testVector(Vector3{1, 1, 0}, UnitX().Add(UnitY()), t)
}

func TestVector3_Subtract(t *testing.T) {
	testVector(Vector3{1, -1, 0}, UnitX().Subtract(UnitY()), t)
}

func TestVector3_MultiplyByConstant(t *testing.T) {
	testVector(Vector3{1.5, 0, 0}, UnitX().MultiplyByConstant(1.5), t)
}

func TestVector3_MultiplyByVector(t *testing.T) {
	testVector(Vector3{1.5, 0, 0}, UnitX().MultiplyByVector(Vector3{1.5, 0, 0}), t)
}

func TestVector3_DivideByConstant(t *testing.T) {
	testVector(Vector3{0.5, 0, 0}, UnitX().DivideByConstant(2.0), t)
}

func TestVector3_DivideByVector(t *testing.T) {
	testVector(Vector3{0.5, 0, 0}, UnitX().DivideByVector(Vector3{2.0, 1.0, 1.0}), t)
}

func TestVector3_Length(t *testing.T) {
	val := Vector3{1, 2, 3}.Length()
	expected := math.Sqrt(14)
	if val != expected {
		t.Errorf("Length expected %f got %f", expected, val)
	}
}

func TestVector3_Normalise(t *testing.T) {
	testVector(Vector3{1, 0, 0}, Vector3{2, 0, 0}.Normalise(), t)
}

func TestVector3_Cross(t *testing.T) {
	testVector(UnitZ(), UnitX().Cross(UnitY()), t)
}

func TestVector3_Dot(t *testing.T) {
	val := Vector3{1, 2, 1}.Dot(Vector3{2, 1, 4})
	expected := 8.0
	if val != expected {
		t.Errorf("Dot product expected %f got %f", expected, val)
	}
}

func TestVector3_Lerp(t *testing.T) {
	testVector(Vector3{0.75, 0.25, 0}, UnitX().Lerp(UnitY(), 0.25), t)
}

func TestVector3_Equal(t *testing.T) {
	testCases := []struct {
		a        Vector3
		b        Vector3
		expected bool
	}{
		{UnitX(), UnitX(), true},
		{UnitX(), UnitY(), false},
		{UnitX(), UnitX().Add(UnitX().MultiplyByConstant(1e-15)), true},
		{UnitX(), UnitX().Add(UnitX().MultiplyByConstant(1e-7)), false},
	}

	for _, testCase := range testCases {
		if result := testCase.a.Equals(testCase.b); result != testCase.expected {
			t.Errorf("%v == %v expected %v, was %v", testCase.a, testCase.b, testCase.expected, result)
		}
	}
}

func TestVector2_Dot(t *testing.T) {
	val := Vector2{1, 2}.Dot(Vector2{2, 4})
	expected := 10.0
	if val != expected {
		t.Errorf("Dot product expected %f got %f", expected, val)
	}
}

func TestVector2_DivideByVector(t *testing.T) {
	testVector2(Vector2{0.5, 0}, Vector2{1.0, 0.0}.DivideByVector(Vector2{2.0, 1.0}), t)
}
