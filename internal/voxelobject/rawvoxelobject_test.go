package voxelobject

import (
	"github.com/mattkimber/gorender/internal/geometry"
	"github.com/mattkimber/gorender/internal/utils/byteutils"
	"testing"
)

func TestRawVoxelObject_Size(t *testing.T) {
	size := geometry.Point{X: 1, Y: 2, Z: 3}
	object := RawVoxelObject(byteutils.Make3DByteSlice(size))
	if object.Size() != size {
		t.Errorf("expected size %v but was %v", size, object.Size())
	}
}
