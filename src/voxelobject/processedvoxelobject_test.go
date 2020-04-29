package voxelobject

import (
	"geometry"
	"testing"
	"utils/fileutils"
	"voxelobject/vox"
)

func TestRawVoxelObject_GetProcessedVoxelObject(t *testing.T) {
	var mv vox.MagicaVoxelObject
	if err := fileutils.InstantiateFromFile("testdata/testcube", &mv); err != nil {
		t.Fatalf("error loading test file: %v", err)
	}

	v := RawVoxelObject(mv).GetProcessedVoxelObject()

	for x := 0; x < len(mv); x++ {
		for y := 0; y < len(mv[x]); y++ {
			for z := 0; z < len(mv[x][y]); z++ {
				if v.SafeGetData(x, y, z).Index != mv[x][y][z] {
					t.Errorf("voxel at [%d,%d,%d] not equal - got %d, expected %d", x, y, z, v.SafeGetData(x, y, z).Index, mv[x][y][z])
				}
			}
		}
	}

	testCases := []struct{
		loc geometry.Point
		expected geometry.Vector3
	}{
		{geometry.Point{},geometry.Vector3{}},
		{geometry.Point{X: 1, Y: 1, Z: 1},geometry.Vector3{X: 1, Y: 1, Z: 1}.Normalise()},
		{geometry.Point{X: 2, Y: 1, Z: 1},geometry.Vector3{X: -1, Y: 1, Z: 1}.Normalise()},
		{geometry.Point{X: 1, Y: 2, Z: 1},geometry.Vector3{X: 1, Y: -1, Z: 1}.Normalise()},
		{geometry.Point{X: 2, Y: 2, Z: 1},geometry.Vector3{X: -1, Y: -1, Z: 1}.Normalise()},
		{geometry.Point{X: 1, Y: 1, Z: 2},geometry.Vector3{X: 1, Y: 1, Z: -1}.Normalise()},
		{geometry.Point{X: 2, Y: 1, Z: 2},geometry.Vector3{X: -1, Y: 1, Z: -1}.Normalise()},
		{geometry.Point{X: 1, Y: 2, Z: 2},geometry.Vector3{X: 1, Y: -1, Z: -1}.Normalise()},
		{geometry.Point{X: 2, Y: 2, Z: 2},geometry.Vector3{X: -1, Y: -1, Z: -1}.Normalise()},
	}

	for _, testCase := range testCases {
		result := v.SafeGetData(testCase.loc.X, testCase.loc.Y, testCase.loc.Z).Normal
		if !result.Equals(testCase.expected) {
			t.Errorf("Normal at %v expected %v, got %v", testCase.loc, testCase.expected, result)
		}
	}
}
