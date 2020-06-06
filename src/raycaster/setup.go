package raycaster

import (
	"geometry"
	"manifest"
	"math"
)

func getRenderDirection(angle float64, elevationAngle float64) geometry.Vector3 {
	x, y, z := -math.Cos(geometry.DegToRad(angle)), math.Sin(geometry.DegToRad(angle)), math.Sin(geometry.DegToRad(elevationAngle))
	return geometry.Vector3{X: x, Y: y, Z: z}.Normalise()
}

func getLightingDirection(angle float64, elevation float64, flipY bool) geometry.Vector3 {
	x, y, z := -math.Cos(geometry.DegToRad(angle)), math.Sin(geometry.DegToRad(angle)), math.Sin(geometry.DegToRad(elevation))
	if flipY {
		y = -y
	}
	return geometry.Zero().Subtract(geometry.Vector3{X: x, Y: y, Z: z}).Normalise()
}

func getViewportPlane(angle float64, m manifest.Manifest, zError float64, size geometry.Point) geometry.Plane {
	elevationAngle := getElevationAngle(m)
	cos, sin := math.Cos(geometry.DegToRad(angle)), math.Sin(geometry.DegToRad(angle))

	midpointX := float64(size.X-1) / 2.0
	if m.PadToFullLength {
		midpointX -= ((m.Size.X - 1) - float64(size.X)) / 2.0
	}

	midpoint := geometry.Vector3{X: midpointX, Y: float64(size.Y) / 2.0, Z: (m.Size.Z + zError) / 2.0}
	viewpoint := midpoint.Add(getRenderDirection(angle, elevationAngle).MultiplyByConstant(m.Size.X / 2))

	planeNormalXComponent := math.Abs(((m.Size.X) / 2.0) * cos * math.Sin(geometry.DegToRad(elevationAngle)))
	planeNormalYComponent := math.Abs(((m.Size.Y) / 2.0) * sin * math.Sin(geometry.DegToRad(elevationAngle)))
	planeNormalZComponent := (m.Size.Z + zError) / 2.0

	planeNormal := geometry.UnitZ().MultiplyByConstant(planeNormalXComponent + planeNormalYComponent + planeNormalZComponent)

	renderNormalXComponent := math.Abs(((m.Size.X) / 2.0) * sin)
	renderNormalYComponent := math.Abs(((m.Size.Y) / 2.0) * cos)
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
	x, y := -math.Cos(geometry.DegToRad(angle)), math.Sin(geometry.DegToRad(angle))
	return geometry.Vector3{X: y, Y: -x}.Normalise()
}
