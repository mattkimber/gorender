package geometry

import gandalfgeo "github.com/mattkimber/gandalf/geometry"

func FromGandalfPoint(point gandalfgeo.Point) Point {
	return Point{
		X: point.X,
		Y: point.Y,
		Z: point.Z,
	}
}

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
