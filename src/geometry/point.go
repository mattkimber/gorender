package geometry

type Point struct {
	X, Y, Z byte
}

type PointWithColour struct {
	Point  Point
	Colour byte
}
