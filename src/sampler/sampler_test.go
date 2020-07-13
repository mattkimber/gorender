package sampler

import (
	"reflect"
	"testing"
)

func TestSquare(t *testing.T) {
	expected := Samples{
		[]Sample{{{}, {Y: 0.5}, {X: 0.25}, {X: 0.25, Y: 0.5}}},
		[]Sample{{{X: 0.5}, {X: 0.5, Y: 0.5}, {X: 0.75}, {X: 0.75, Y: 0.5}}},
	}

	if gotResult := Square(2, 1, 2, 0); !reflect.DeepEqual(gotResult, expected) {
		t.Errorf("Square() = %v, want %v", gotResult, expected)
	}
}

func TestDisc(t *testing.T) {
	result := Disc(5, 5, 3, .1)
	if len(result[0][0]) != 9 {
		t.Errorf("Disc() = %d, want %d", len(result[0][0]), 9)
	}
}
