package raycaster

import (
	"geometry"
	"math"
)

func getRenderDirection(angle int) geometry.Vector3 {
	x, y, z := -math.Cos(degToRad(angle)), math.Sin(degToRad(angle)), math.Sin(degToRad(30))
	return geometry.Vector3{X: x, Y: y, Z: z}.Normalise()
}

func degToRad(angle int) float64 {
	return (float64(angle) / 180.0) * math.Pi
}

func getLightingDirection(angle int) geometry.Vector3 {
	x, y, z := -math.Cos(degToRad(angle)), math.Sin(degToRad(angle)), math.Sin(degToRad(lightingElevationAngle))
	return geometry.Zero().Subtract(geometry.Vector3{X: x, Y: y, Z: z}).Normalise()
}

func getViewportPlane(angle int, size geometry.Point) geometry.Plane {
	midpoint := geometry.Vector3{X: float64(size.X) / 2.0, Y: float64(size.Y) / 2.0, Z: float64(size.Y) / 2.0}
	viewpoint := midpoint.Add(getRenderDirection(angle).MultiplyByConstant(100.0))

	planeNormal := geometry.UnitZ().MultiplyByConstant(midpoint.X)
	renderNormal := getRenderNormal(angle).MultiplyByConstant(midpoint.X)

	a := viewpoint.Subtract(renderNormal).Subtract(planeNormal)
	b := viewpoint.Add(renderNormal).Subtract(planeNormal)
	c := viewpoint.Add(renderNormal).Add(planeNormal)
	d := viewpoint.Subtract(renderNormal).Add(planeNormal)

	return geometry.Plane{A: a, B: b, C: c, D: d}
}

func getRenderNormal(angle int) geometry.Vector3 {
	x, y := -math.Cos(degToRad(angle)), math.Sin(degToRad(angle))
	return geometry.Vector3{X: y, Y: -x}.Normalise()
}
