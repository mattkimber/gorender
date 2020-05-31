package sampler

import (
	"reflect"
	"testing"
)

func TestSquare(t *testing.T) {
	expected := Samples{
		[]Sample{{{0, 0}, {0, 0.5}, {0.25, 0}, {0.25, 0.5}}},
		[]Sample{{{0.5, 0}, {0.5, 0.5}, {0.75, 0}, {0.75, 0.5}}},
	}

	if gotResult := Square(2, 1, 2, 0); !reflect.DeepEqual(gotResult, expected) {
		t.Errorf("Square() = %v, want %v", gotResult, expected)
	}
}
