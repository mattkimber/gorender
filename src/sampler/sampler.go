package sampler

import "geometry"

type Sample []geometry.Vector2

type Samples [][]Sample

func (s Samples) Width() int {
	return len(s)
}

func (s Samples) Height() int {
	return len(s[0])
}

func Square(width, height int, scale int) (result Samples) {

	result = make([][]Sample, width)
	for i := 0; i < width; i++ {
		result[i] = make([]Sample, height)
		for j := 0; j < height; j++ {
			result[i][j] = make(Sample, scale*scale)
			for k := 0; k < scale; k++ {
				for l := 0; l < scale; l++ {
					result[i][j][l+(k*scale)] = geometry.Vector2{
						X: (float64(i*scale) + float64(k)) / (float64(width * scale)),  //+ (float64(k) / float64(width * scale)),
						Y: (float64(j*scale) + float64(l)) / (float64(height * scale)), //+ (float64(l) / float64(height * scale)),
					}
				}
			}
		}
	}

	return result
}
