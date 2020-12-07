package sampler

import (
	"github.com/mattkimber/gorender/internal/geometry"
	"reflect"
	"testing"
)

func TestSquare(t *testing.T) {
	expected := Samples{
		[]SampleList{{
			{ Location:  geometry.Vector2{X: 1.0/6.0, Y: 1.0/3.0}, Influence: 0.9444444444444444},
			{ Location:  geometry.Vector2{X: 1.0/6.0, Y: 2.0/3.0}, Influence: 0.9444444444444444},
			{ Location:  geometry.Vector2{X: 2.0/6.0, Y: 1.0/3.0}, Influence: 0.9444444444444444},
			{ Location:  geometry.Vector2{X: 2.0/6.0, Y: 2.0/3.0}, Influence: 0.9444444444444444},
		}},
		[]SampleList{{
			{ Location:  geometry.Vector2{X: 4.0/6.0, Y: 1.0/3.0}, Influence: 0.9444444444444444},
			{ Location:  geometry.Vector2{X: 4.0/6.0, Y: 2.0/3.0}, Influence: 0.9444444444444444},
			{ Location:  geometry.Vector2{X: 0.8333333333333333, Y: 1.0/3.0}, Influence: 0.9444444444444444},
			{ Location:  geometry.Vector2{X: 0.8333333333333333, Y: 2.0/3.0}, Influence: 0.9444444444444444},
		}},
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
