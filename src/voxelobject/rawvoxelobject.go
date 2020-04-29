package voxelobject

import (
	"geometry"
)

type RawVoxelObject [][][]byte

func (v RawVoxelObject) Size() geometry.Point {
	if v == nil || len(v) == 0 || len(v[0]) == 0 {
		return geometry.Point{}
	}

	return geometry.Point{
		X: len(v),
		Y: len(v[0]),
		Z: len(v[0][0]),
	}
}

func (v RawVoxelObject) Invalid() bool {
	size := v.Size()
	return size.X == 0 || size.Y == 0 || size.Z == 0
}
