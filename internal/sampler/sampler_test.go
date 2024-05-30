package sampler

import (
	"github.com/mattkimber/gorender/internal/geometry"
	"math"
	"testing"
)

func TestSquare(t *testing.T) {
	expected := Samples{
		[]SampleList{{
			{Location: geometry.Vector2{X: 1.0 / 6.0, Y: 1.0 / 3.0}, Influence: 0.5285954792089682},
			{Location: geometry.Vector2{X: 1.0 / 6.0, Y: 2.0 / 3.0}, Influence: 0.5285954792089683},
			{Location: geometry.Vector2{X: 2.0 / 6.0, Y: 1.0 / 3.0}, Influence: 0.5285954792089683},
			{Location: geometry.Vector2{X: 2.0 / 6.0, Y: 2.0 / 3.0}, Influence: 0.5285954792089684},
		}},
		[]SampleList{{
			{Location: geometry.Vector2{X: 4.0 / 6.0, Y: 1.0 / 3.0}, Influence: 0.5285954792089682},
			{Location: geometry.Vector2{X: 4.0 / 6.0, Y: 2.0 / 3.0}, Influence: 0.5285954792089683},
			{Location: geometry.Vector2{X: 0.8333333333333333, Y: 1.0 / 3.0}, Influence: 0.5285954792089683},
			{Location: geometry.Vector2{X: 0.8333333333333333, Y: 2.0 / 3.0}, Influence: 0.5285954792089684},
		}},
	}

	// This test is infuriatingly clunky due to Mac and Win/Linux returning different floating point
	// roundings in GitHub actions
	res := Square(2, 1, 2, 0, 0.5)
	if len(res) != len(expected) {
		t.Errorf("Square() arrays not even, got %d expected %d", len(res), len(expected))
	}

	for i, _ := range res {
		if len(res[i]) != len(expected[i]) {
			t.Errorf("Square() %d arrays not even, got %d expected %d", i, len(res[i]), len(expected[i]))
		}

		for j, _ := range res[i] {
			if len(res[j]) != len(expected[j]) {
				t.Errorf("Square() %d.%d arrays not even, got %d expected %d", i, j, len(res[j]), len(expected[j]))
			}

			for k, _ := range res[j] {
				if math.Abs(expected[i][j][k].Location.X-res[i][j][k].Location.X) > 0.00001 {
					t.Errorf("Square() %d.%d.%d X = %f, want %f", i, j, k, res[i][j][k].Location.X, expected[i][j][k].Location.X)
				}

				if math.Abs(expected[i][j][k].Location.Y-res[i][j][k].Location.Y) > 0.00001 {
					t.Errorf("Square() %d.%d.%d Y = %f, want %f", i, j, k, res[i][j][k].Location.Y, expected[i][j][k].Location.Y)
				}

				if math.Abs(expected[i][j][k].Influence-res[i][j][k].Influence) > 0.00001 {
					t.Errorf("Square() %d.%d.%d Z = %f, want %f", i, j, k, res[i][j][k].Influence, expected[i][j][k].Influence)
				}
			}
		}
	}
}

func TestDisc(t *testing.T) {
	result := Disc(5, 5, 3, .1, 0.5)
	if len(result[0][0]) != 9 {
		t.Errorf("Disc() = %d, want %d", len(result[0][0]), 9)
	}
}
