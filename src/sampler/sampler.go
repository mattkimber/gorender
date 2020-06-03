package sampler

import (
	"geometry"
	"math/rand"
)

type Sample []geometry.Vector2

type Samples [][]Sample

func (s Samples) Width() int {
	return len(s)
}

func (s Samples) Height() int {
	return len(s[0])
}

func Get(name string) func(int, int, int, float64) Samples {
	switch name {
	case "square":
		return Square
	case "disc":
		return Disc
	default:
		return Square
	}
}

func Square(width, height int, accuracy int, overlap float64) (result Samples) {
	result = make([][]Sample, width)
	for i := 0; i < width; i++ {
		result[i] = make([]Sample, height)
		for j := 0; j < height; j++ {
			result[i][j] = make(Sample, accuracy*accuracy)
			for k := 0; k < accuracy; k++ {
				for l := 0; l < accuracy; l++ {
					result[i][j][l+(k*accuracy)] = geometry.Vector2{
						X: (float64(i*accuracy) + (float64(k) * (1.0 + overlap))) / (float64(width * accuracy)),
						Y: (float64(j*accuracy) + (float64(l) * (1.0 + overlap))) / (float64(height * accuracy)),
					}
				}
			}
		}
	}

	return result
}

const discs = 10

var discCache [][]geometry.Vector2

func Disc(width, height int, accuracy int, overlap float64) (result Samples) {
	result = make([][]Sample, width)
	scaleVec := geometry.Vector2{X: float64(width), Y: float64(height)}
	for i := 0; i < width; i++ {
		result[i] = make([]Sample, height)
		for j := 0; j < height; j++ {
			loc := geometry.Vector2{X: float64(i) / scaleVec.X, Y: float64(j) / scaleVec.Y}
			disc := getPoissonDisc(accuracy, overlap)

			result[i][j] = make(Sample, len(disc))
			for k, s := range disc {
				result[i][j][k] = loc.Add(s.DivideByVector(scaleVec))
			}
		}
	}

	return
}

// Get a poisson disc using the naive/slow dart throwing algorithm
func getPoissonDisc(accuracy int, overlap float64) []geometry.Vector2 {
	if discCache == nil {
		discCache = make([][]geometry.Vector2, discs)
	}

	discNum := rand.Intn(discs)
	if discCache[discNum] != nil {
		return discCache[discNum]
	}

	numSamples := accuracy * accuracy
	distance := 1.0 / float64(accuracy)
	distance = distance * distance

	radius := 1.0 + overlap

	disc := make([]geometry.Vector2, 0)
	var valid bool

	// Create a poisson disc by dart throwing
	for i := 0; i < numSamples*1000; i++ {
		valid = true
		trial := geometry.Vector2{X: (rand.Float64() - 0.5) * radius, Y: (rand.Float64() - 0.5) * radius}
		for k := 0; k < len(disc); k++ {
			if trial.LengthSquared() > radius || trial.DistanceSquared(disc[k]) < distance {
				valid = false
			}
		}

		if valid {
			disc = append(disc, trial)
			if len(disc) >= numSamples {
				break
			}
		}
	}

	discCache[discNum] = disc
	return disc
}
