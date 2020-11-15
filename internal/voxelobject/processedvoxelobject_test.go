package voxelobject

import (
	"github.com/mattkimber/gorender/internal/colour"
	"github.com/mattkimber/gorender/internal/geometry"
	"github.com/mattkimber/gorender/internal/utils/fileutils"
	"github.com/mattkimber/gorender/internal/voxelobject/vox"
	"testing"
)

func TestRawVoxelObject_GetProcessedVoxelObject(t *testing.T) {
	var mv vox.MagicaVoxelObject
	if err := fileutils.InstantiateFromFile("testdata/testcube", &mv); err != nil {
		t.Fatalf("error loading test file: %v", err)
	}

	v := RawVoxelObject(mv).GetProcessedVoxelObject(&colour.Palette{}, false, false)
	testObject(t, mv, v)

	v = RawVoxelObject(mv).GetProcessedVoxelObject(&colour.Palette{}, true, false)
	testObject(t, mv, v)

	v = RawVoxelObject(mv).GetProcessedVoxelObject(&colour.Palette{}, true, true)
	testObject(t, mv, v)

	v = RawVoxelObject(mv).GetProcessedVoxelObject(&colour.Palette{}, false, true)
	testObject(t, mv, v)
}

func testObject(t *testing.T, mv vox.MagicaVoxelObject, v ProcessedVoxelObject) {
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
		{geometry.Point{X: 3, Y: 2, Z: 2}, geometry.Vector3{X: 0.341305127505155, Y: 0.6646468272468807, Z: 0.6646468272468807}},
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

	v := RawVoxelObject(mv).GetProcessedVoxelObject(&colour.Palette{}, false, false)
	return v
}

func TestProcessedVoxelObject_getOcclusion(t *testing.T) {
	p := getObject("occlude", t)
	if p.getOcclusion(2, 2, 2) != 1 {
		t.Errorf("Occlusion at 2,2,2 is %d, expected 1\n", p.getOcclusion(2, 2, 2))
	}

	if p.getOcclusion(3, 3, 3) != 0 {
		t.Errorf("Occlusion at 3,3,3 is %d, expected 0\n", p.getOcclusion(3, 3, 3))
	}

}
