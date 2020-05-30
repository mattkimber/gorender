package geometry

type Point struct {
	X, Y, Z int
}

type PointWithColour struct {
	Point  Point
	Colour byte
}

func (p *Point) ToVector3() Vector3 {
	return Vector3{
		X: float64(p.X),
		Y: float64(p.Y),
		Z: float64(p.Z),
	}
}