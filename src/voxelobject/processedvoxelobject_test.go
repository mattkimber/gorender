package voxelobject

import (
	"colour"
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

	v := RawVoxelObject(mv).GetProcessedVoxelObject(&colour.Palette{})

	for x := 0; x < len(mv); x++ {
		for y := 0; y < len(mv[x]); y++ {
			for z := 0; z < len(mv[x][y]); z++ {
				if v.SafeGetData(x, y, z).Index != mv[x][y][z] {
					t.Errorf("voxel at [%d,%d,%d] not equal - got %d, expected %d", x, y, z, v.SafeGetData(x, y, z).Index, mv[x][y][z])
				}
			}
		}
	}
}

func TestRawVoxelObject_GetProcessedVoxelObject_Normals(t *testing.T) {
	v := getObject("testcube", t)

	testCases := []struct {
		loc      geometry.Point
		expected geometry.Vector3
	}{
		{geometry.Point{}, geometry.Vector3{}},
		{geometry.Point{X: 1, Y: 1, Z: 1}, geometry.Vector3{X: 1, Y: 1, Z: 1}.Normalise()},
		{geometry.Point{X: 2, Y: 1, Z: 1}, geometry.Vector3{X: -1, Y: 1, Z: 1}.Normalise()},
		{geometry.Point{X: 1, Y: 2, Z: 1}, geometry.Vector3{X: 1, Y: -1, Z: 1}.Normalise()},
		{geometry.Point{X: 2, Y: 2, Z: 1}, geometry.Vector3{X: -1, Y: -1, Z: 1}.Normalise()},
		{geometry.Point{X: 1, Y: 1, Z: 2}, geometry.Vector3{X: 1, Y: 1, Z: -1}.Normalise()},
		{geometry.Point{X: 2, Y: 1, Z: 2}, geometry.Vector3{X: -1, Y: 1, Z: -1}.Normalise()},
		{geometry.Point{X: 1, Y: 2, Z: 2}, geometry.Vector3{X: 1, Y: -1, Z: -1}.Normalise()},
		{geometry.Point{X: 2, Y: 2, Z: 2}, geometry.Vector3{X: -1, Y: -1, Z: -1}.Normalise()},
	}

	for _, testCase := range testCases {
		result := v.SafeGetData(testCase.loc.X, testCase.loc.Y, testCase.loc.Z).Normal
		if !result.Equals(testCase.expected) {
			t.Errorf("Normal at %v expected %v, got %v", testCase.loc, testCase.expected, result)
		}
	}
}

func TestRawVoxelObject_GetProcessedVoxelObject_AveragedNormals(t *testing.T) {
	v := getObject("testcube_big", t)

	testCases := []struct {
		loc      geometry.Point
		expected geometry.Vector3
	}{
		{geometry.Point{}, geometry.Vector3{}},
		{geometry.Point{X: 3, Y: 2, Z: 2}, geometry.Vector3{X: 0.43104479944651053, Y: 0.6403956621583834, Z: 0.635683708111859}},
		{geometry.Point{X: 9, Y: 2, Z: 2}, geometry.Vector3{X: -1, Y: 1, Z: 1}.Normalise()},
		{geometry.Point{X: 2, Y: 9, Z: 2}, geometry.Vector3{X: 1, Y: -1, Z: 1}.Normalise()},
		{geometry.Point{X: 9, Y: 9, Z: 2}, geometry.Vector3{X: -1, Y: -1, Z: 1}.Normalise()},
		{geometry.Point{X: 2, Y: 2, Z: 9}, geometry.Vector3{X: 0.5773502691896258, Y: 0.5773502691896258, Z: -0.5773502691896258}},
		{geometry.Point{X: 9, Y: 2, Z: 9}, geometry.Vector3{X: -1, Y: 1, Z: -1}.Normalise()},
		{geometry.Point{X: 2, Y: 9, Z: 9}, geometry.Vector3{X: 1, Y: -1, Z: -1}.Normalise()},
		{geometry.Point{X: 9, Y: 9, Z: 9}, geometry.Vector3{X: -1, Y: -1, Z: -1}.Normalise()},
	}

	for _, testCase := range testCases {
		result := v.SafeGetData(testCase.loc.X, testCase.loc.Y, testCase.loc.Z).AveragedNormal
		if !result.Equals(testCase.expected) {
			t.Errorf("Average normal at %v expected %v, got %v", testCase.loc, testCase.expected, result)
		}
	}
}

func getObject(filename string, t *testing.T) ProcessedVoxelObject {
	var mv vox.MagicaVoxelObject
	if err := fileutils.InstantiateFromFile("testdata/"+filename, &mv); err != nil {
		t.Fatalf("error loading test file: %v", err)
	}

	v := RawVoxelObject(mv).GetProcessedVoxelObject(&colour.Palette{})
	return v
}
