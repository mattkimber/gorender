package geometry

import "math"

type Plane struct {
	A, B, C, D Vector3
}

func (a Plane) Equals(b Plane) bool {
	return a.A.Equals(b.A) && a.B.Equals(b.B) && a.C.Equals(b.C) && a.D.Equals(b.D)
}

func (a Plane) BiLerpWithinPlane(u float64, v float64) Vector3 {
	abu, dcu := a.A.Lerp(a.B, u), a.D.Lerp(a.C, u)
	return dcu.Lerp(abu, v)
}

func DegToRad(angle float64) float64 {
	return (angle / 180.0) * math.Pi
}
