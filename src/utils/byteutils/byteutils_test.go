package byteutils

import (
	"geometry"
	"testing"
)

func TestMake3DByteSlice(t *testing.T) {
	slice := Make3DByteSlice(geometry.Point{X: 1, Y: 2, Z: 3})

	if len(slice) != 1 {
		t.Fatalf("x length is incorrect: expected 1, was %d", len(slice))
	}

	if len(slice[0]) != 2 {
		t.Fatalf("y length is incorrect: expected 2, was %d", len(slice[0]))
	}

	if len(slice[0][0]) != 3 {
		t.Fatalf("y length is incorrect: expected 2, was %d", len(slice[0][0]))
	}

}
