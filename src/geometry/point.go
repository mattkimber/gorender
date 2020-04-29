package geometry

type Point struct {
	X, Y, Z int
}

type PointWithColour struct {
	Point  Point
	Colour byte
}
