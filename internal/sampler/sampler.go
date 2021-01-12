package sampler

import (
	"github.com/mattkimber/gorender/internal/geometry"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
)

type Sample struct {
	Location geometry.Vector2
	Influence float64
}

type SampleList []Sample

type Samples [][]SampleList

func (s Samples) Width() int {
	return len(s)
}

func (s Samples) Height() int {
	return len(s[0])
}

func (s Samples) GetImage() (img *image.RGBA) {
	rect := image.Rect(0, 0, 200, 200)
	img = image.NewRGBA(rect)

	// Clear to white
	draw.Draw(img, rect, image.NewUniform(color.White), image.Point{}, draw.Over)

	samples := s[0][0]
	for _, smp := range samples {
		x, y := int(100.0+(smp.Location.X*50.0)), int(100.0+(smp.Location.Y*50.0))
		if x >= 0 && y >= 0 && x < 200 && y < 200 {
			img.Set(x, y, color.RGBA{R: uint8(smp.Influence*255.0), G: 0, B: 0, A: 255})
		}
	}

	return
}

func Get(name string) func(int, int, int, float64, float64) Samples {
	switch name {
	case "square":
		return Square
	case "disc":
		return Disc
	default:
		return Square
	}
}

func Square(width, height int, accuracy int, overlap float64, falloff float64) (result Samples) {
	fAccuracy := float64(accuracy)

	centre := geometry.Vector2{
		X: 0.5,
		Y: 0.5,
	}

	var location geometry.Vector2

	result = make([][]SampleList, width)
	for i := 0; i < width; i++ {
		result[i] = make([]SampleList, height)
		for j := 0; j < height; j++ {
			result[i][j] = make(SampleList, accuracy*accuracy)

			for k := 0; k < accuracy; k++ {
				fractionK := (1.0 + float64(k)) / (1.0 + fAccuracy)

				for l := 0; l < accuracy; l++ {
					fractionL := (1.0 + float64(l)) / (1.0 + fAccuracy)

					fraction := geometry.Vector2{
						X: fractionK,
						Y: fractionL,
					}


					location = geometry.Vector2{
						X: (float64(i*accuracy) + (fractionK * (1.0 + overlap))*fAccuracy) / (float64(width*accuracy)),
						Y: (float64(j*accuracy) + (fractionL * (1.0 + overlap))*fAccuracy) / (float64(height*accuracy)),
					}

					influence := 1.0 - (math.Pow(centre.DistanceSquared(fraction), falloff) * 2.0)

					if influence < 0 {
						influence = 0
					}

					result[i][j][l+(k*accuracy)] = Sample{
						Location: location,
						Influence: influence,
					}
				}
			}
		}
	}

	return result
}

const discs = 10

var discCache [][]geometry.Vector2

func Disc(width, height int, accuracy int, overlap float64, falloff float64) (result Samples) {
	radiusSquared := (0.5 + overlap) * (0.5 + overlap)
	var location geometry.Vector2

	result = make([][]SampleList, width)
	scaleVec := geometry.Vector2{X: float64(width), Y: float64(height)}
	for i := 0; i < width; i++ {
		result[i] = make([]SampleList, height)
		for j := 0; j < height; j++ {
			loc := geometry.Vector2{X: float64(i) / scaleVec.X, Y: float64(j) / scaleVec.Y}
			disc := getPoissonDisc(accuracy, overlap)

			result[i][j] = make(SampleList, len(disc))
			for k, s := range disc {
				location = loc.Add(s.DivideByVector(scaleVec))

				influence := 1.0 - (math.Pow(radiusSquared, falloff))

				if influence < 0 {
					influence = 0
				}

				result[i][j][k] = Sample{
					Location: location,
					Influence: influence,
				}
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

	radius := 0.5 + overlap

	disc := make([]geometry.Vector2, 0)
	var valid bool

	// Create a poisson disc by dart throwing
	for i := 0; i < numSamples*1000; i++ {
		valid = true
		trial := geometry.Vector2{X: (rand.Float64() - 0.5) * 2.0 * radius, Y: (rand.Float64() - 0.5) * 2.0 * radius}
		for k := 0; k < len(disc); k++ {
			if trial.LengthSquared() > radius*radius || trial.DistanceSquared(disc[k]) < distance {
				valid = false
				break
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
