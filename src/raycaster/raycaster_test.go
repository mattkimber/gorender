package raycaster

import (
	"colour"
	"geometry"
	"testing"
	"utils/fileutils"
	"voxelobject"
	"voxelobject/vox"
)

func Test_getLightingValue(t *testing.T) {
	testCases := []struct {
		normal, lighting geometry.Vector3
		expected         float64
	}{
		{geometry.Vector3{}, geometry.UnitX(), 0.0},
		{geometry.UnitX(), geometry.UnitX(), 1.0},
		{geometry.Vector3{X: 0.5, Y: 1}.Normalise(), geometry.Vector3{X: 1, Y: 0.5, Z: 1}.Normalise(), 0.5962847939999438},
	}
	for _, testCase := range testCases {
		if result := getLightingValue(testCase.normal, testCase.lighting); result != testCase.expected {
			t.Errorf("getLightingValue for normal %v and lighting %v returned %v, expected %v", testCase.normal, testCase.lighting, result, testCase.expected)
		}
	}
}

func getObject(filename string, t *testing.T) voxelobject.ProcessedVoxelObject {
	var mv vox.MagicaVoxelObject
	if err := fileutils.InstantiateFromFile("testdata/"+filename, &mv); err != nil {
		t.Fatalf("error loading test file: %v", err)
	}

	v := voxelobject.RawVoxelObject(mv).GetProcessedVoxelObject(&colour.Palette{})
	return v
}
