package byteutils

import (
	"github.com/mattkimber/gorender/internal/geometry"
)

func Make3DByteSlice(size geometry.Point) [][][]byte {
	result := make([][][]byte, size.X)

	for x := range result {
		result[x] = make([][]byte, size.Y)
		for y := range result[x] {
			result[x][y] = make([]byte, size.Z)
		}
	}

	return result
}
