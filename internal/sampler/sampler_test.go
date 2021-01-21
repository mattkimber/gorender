package sampler

import (
	"github.com/mattkimber/gorender/internal/geometry"
	"reflect"
	"testing"
)

func TestSquare(t *testing.T) {
	expected := Samples{
		[]SampleList{{
			{ Location:  geometry.Vector2{X: 1.0/6.0, Y: 1.0/3.0}, Influence: 0.5285954792089682},
			{ Location:  geometry.Vector2{X: 1.0/6.0, Y: 2.0/3.0}, Influence: 0.5285954792089683},
			{ Location:  geometry.Vector2{X: 2.0/6.0, Y: 1.0/3.0}, Influence: 0.5285954792089683},
			{ Location:  geometry.Vector2{X: 2.0/6.0, Y: 2.0/3.0}, Influence: 0.5285954792089684},
		}},
		[]SampleList{{
			{ Location:  geometry.Vector2{X: 4.0/6.0, Y: 1.0/3.0}, Influence: 0.5285954792089682},
			{ Location:  geometry.Vector2{X: 4.0/6.0, Y: 2.0/3.0}, Influence: 0.5285954792089683},
			{ Location:  geometry.Vector2{X: 0.8333333333333333, Y: 1.0/3.0}, Influence: 0.5285954792089683},
			{ Location:  geometry.Vector2{X: 0.8333333333333333, Y: 2.0/3.0}, Influence: 0.5285954792089684},
		}},
	}

	if gotResult := Square(2, 1, 2, 0, 0.5); !reflect.DeepEqual(gotResult, expected) {
		t.Errorf("Square() = %v, want %v", gotResult, expected)
	}
}

func TestDisc(t *testing.T) {
	result := Disc(5, 5, 3, .1, 0.5)
	if len(result[0][0]) != 9 {
		t.Errorf("Disc() = %d, want %d", len(result[0][0]), 9)
	}
}
