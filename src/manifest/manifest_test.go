package manifest

import (
	"geometry"
	"os"
	"reflect"
	"testing"
)

func TestFromJson(t *testing.T) {
	expected := Manifest{
		LightingAngle:     60,
		LightingElevation: 65,
		DepthInfluence:    0.2,
		Size: geometry.Point{
			X: 20,
			Y: 30,
			Z: 40,
		},
	}

	file, err := os.Open("testdata/manifest.json")
	if err != nil {
		t.Fatalf("Could not open test data: %v", err)
	}

	defer file.Close()

	actual, err := FromJson(file)
	if err != nil {
		t.Fatalf("Could not process test data: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
