package voxelobject

type Point struct {
	X, Y, Z byte
}

type PointWithColour struct {
	Point  Point
	Colour byte
}

type RawVoxelObject [][][]byte

func MakeRawVoxelObject(size Point) RawVoxelObject {
	result := make([][][]byte, size.X)

	for x := range result {
		result[x] = make([][]byte, size.Y)
		for y := range result[x] {
			result[x][y] = make([]byte, size.Z)
		}
	}

	return result
}

func (o RawVoxelObject) Size() Point {
	if o == nil || len(o) == 0 || len(o[0]) == 0 {
		return Point{}
	}

	return Point{
		X: byte(len(o)),
		Y: byte(len(o[0])),
		Z: byte(len(o[0][0])),
	}
}

func (o RawVoxelObject) Invalid() bool {
	size := o.Size()
	return size.X == 0 || size.Y == 0 || size.Z == 0
}
