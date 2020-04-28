package voxelobject

import "testing"

func TestRawVoxelObject_Size(t *testing.T) {
	size := Point{X: 1, Y: 2, Z: 3}
	object := MakeRawVoxelObject(size)
	if object.Size() != size {
		t.Errorf("expected size %v but was %v", size, object.Size())
	}
}
