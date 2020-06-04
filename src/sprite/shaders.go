package sprite

import (
	"colour"
	"geometry"
	"manifest"
	"raycaster"
)

func Colour(smp raycaster.RenderSample, d manifest.Definition, resolveSpecialColours bool) colour.RGB {
	lightingOffset := getLightingOffset(smp, d.Manifest.DepthInfluence)
	return d.Palette.GetLitRGB(smp.Index, lightingOffset, d.Manifest.Brightness, d.Manifest.Contrast, resolveSpecialColours)
}

func Normal(smp raycaster.RenderSample) colour.RGB {
	normal := smp.Normal.MultiplyByConstant(32766).Add(geometry.Vector3{X: 32766, Y: 32766, Z: 32766})
	return colour.RGB{R: normal.X, G: normal.Y, B: normal.Z}
}

func AveragedNormal(smp raycaster.RenderSample) colour.RGB {
	normal := smp.AveragedNormal.MultiplyByConstant(32766).Add(geometry.Vector3{X: 32766, Y: 32766, Z: 32766})
	return colour.RGB{R: normal.X, G: normal.Y, B: normal.Z}
}

func Depth(smp raycaster.RenderSample) colour.RGB {
	v := float64(smp.Depth * 400)
	return colour.RGB{R: v, G: v, B: v}
}

func Occlusion(smp raycaster.RenderSample) colour.RGB {
	v := float64(smp.Occlusion * 6000)
	return colour.RGB{R: v, G: v, B: v}
}

func Shadow(smp raycaster.RenderSample) colour.RGB {
	v := 65535 - (smp.Shadowing * 65535)
	return colour.RGB{R: v, G: v, B: v}
}

func Lighting(smp raycaster.RenderSample) colour.RGB {
	v := 32767 + (smp.LightAmount * 32767)
	return colour.RGB{R: v, G: v, B: v}
}

/*
func Apply32bppSprite(img *image.RGBA, bounds image.Rectangle, loc image.Point, info raycaster.RenderOutput, d manifest.Definition) {
	shader := func(smp raycaster.RenderSample) (float64, float64, float64) {
		lightingOffset := getLightingOffset(smp, d.Manifest.DepthInfluence)
		return d.Palette.GetLitRGB(smp.Index, lightingOffset, d.Manifest.Brightness, d.Manifest.Contrast)
	}

	apply32bppImage(img, bounds, loc, shader, info, d.SoftenEdges())
}
*/
/*
func ApplyIndexedSprite(img *image.Paletted, bounds image.Rectangle, loc image.Point, info raycaster.RenderOutput, d manifest.Definition) {
	shader := func(smp raycaster.RenderSample) byte {
		lightingOffset := getLightingOffset(smp, d.Manifest.DepthInfluence)
		idx := d.Palette.GetLitIndexed(smp.Index, lightingOffset)
		return idx
	}

	applyIndexedImage(img, d.Palette, bounds, loc, shader, info)
}

func ApplyMaskSprite(img *image.Paletted, bounds image.Rectangle, loc image.Point, info raycaster.RenderOutput, d manifest.Definition) {
	shader := func(smp raycaster.RenderSample) byte {
		return d.Palette.GetMaskColour(smp.Index)
	}

	applyIndexedImage(img, d.Palette, bounds, loc, shader, info)
}
*/

func getLightingOffset(smp raycaster.RenderSample, depthInfluence float64) float64 {
	lightingOffset := -0.3
	lightingOffset += smp.LightAmount * 0.6
	lightingOffset += (-(float64(smp.Depth-120) / 40)) * depthInfluence
	lightingOffset += (-float64(smp.Occlusion) / 10.0) * 0.2
	lightingOffset -= smp.Shadowing * 0.2

	lightingOffset = lightingOffset / 1.5

	return lightingOffset
}
