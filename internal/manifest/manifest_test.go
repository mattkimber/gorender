package manifest

import (
	"github.com/mattkimber/gorender/internal/geometry"
	"os"
	"reflect"
	"testing"
)

func TestFromJson(t *testing.T) {
	expected := Manifest{
		LightingAngle:     60,
		LightingElevation: 65,
		DepthInfluence:    0.2,
		Accuracy:          2,
		Contrast:          1.0,
		EdgeThreshold:     0.5,
		TilingMode:        "normal",
		Size: geometry.Vector3{
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
