package raycaster

import (
	"geometry"
	"manifest"
	"math"
)

func getRenderDirection(angle float64, elevationAngle float64) geometry.Vector3 {
	x, y, z := -math.Cos(degToRad(angle)), math.Sin(degToRad(angle)), math.Sin(degToRad(elevationAngle))
	return geometry.Vector3{X: x, Y: y, Z: z}.Normalise()
}

func degToRad(angle float64) float64 {
	return (angle / 180.0) * math.Pi
}

func getLightingDirection(angle float64, elevation float64) geometry.Vector3 {
	x, y, z := -math.Cos(degToRad(angle)), math.Sin(degToRad(angle)), math.Sin(degToRad(elevation))
	return geometry.Zero().Subtract(geometry.Vector3{X: x, Y: y, Z: z}).Normalise()
}

func getViewportPlane(angle float64, m manifest.Manifest, size geometry.Point) geometry.Plane {
	elevationAngle := getElevationAngle(m)
	cos, sin := math.Cos(degToRad(angle)), math.Sin(degToRad(angle))

	midpoint := geometry.Vector3{X: float64(size.X) / 2.0, Y: float64(size.Y) / 2.0, Z: float64(m.Size.Z) / 2.0}
	viewpoint := midpoint.Add(getRenderDirection(angle, elevationAngle).MultiplyByConstant(100.0))

	planeNormalXComponent := math.Abs((float64(m.Size.X) / 2.0) * cos * math.Sin(degToRad(elevationAngle)))
	planeNormalYComponent := math.Abs((float64(m.Size.Y) / 2.0) * sin * math.Sin(degToRad(elevationAngle)))
	planeNormalZComponent := float64(m.Size.Z) / 2.0
	planeNormal := geometry.UnitZ().MultiplyByConstant(planeNormalXComponent + planeNormalYComponent + planeNormalZComponent)

	renderNormalXComponent := math.Abs((float64(m.Size.X) / 2.0) * sin)
	renderNormalYComponent := math.Abs((float64(m.Size.Y) / 2.0) * cos)
	renderNormal := getRenderNormal(angle).MultiplyByConstant(renderNormalXComponent + renderNormalYComponent)

	a := viewpoint.Subtract(renderNormal).Subtract(planeNormal)
	b := viewpoint.Add(renderNormal).Subtract(planeNormal)
	c := viewpoint.Add(renderNormal).Add(planeNormal)
	d := viewpoint.Subtract(renderNormal).Add(planeNormal)

	return geometry.Plane{A: a, B: b, C: c, D: d}
}

func getElevationAngle(m manifest.Manifest) float64 {
	return float64(m.RenderElevationAngle)
}

func getRenderNormal(angle float64) geometry.Vector3 {
	x, y := -math.Cos(degToRad(angle)), math.Sin(degToRad(angle))
	return geometry.Vector3{X: y, Y: -x}.Normalise()
}
