package spritesheet

import (
	"github.com/mattkimber/gandalf/magica"
	"github.com/mattkimber/gorender/internal/colour"
	"github.com/mattkimber/gorender/internal/geometry"
	"github.com/mattkimber/gorender/internal/manifest"
	"github.com/mattkimber/gorender/internal/raycaster"
	"github.com/mattkimber/gorender/internal/sampler"
	"github.com/mattkimber/gorender/internal/sprite"
	"github.com/mattkimber/gorender/internal/utils/imageutils"
	"github.com/mattkimber/gorender/internal/voxelobject"
	"image"
	"os"
	"testing"
)

func TestGetSpritesheets(t *testing.T) {
	def := manifest.Definition{
		Palette: colour.Palette{Entries: []colour.PaletteEntry{{R: 0, G: 0, B: 0}, {R: 255, G: 255, B: 255}}},
		Scale:   1.0,
		Manifest: manifest.Manifest{
			LightingAngle:        45,
			LightingElevation:    60,
			Size:                 geometry.Vector3{},
			RenderElevationAngle: 0,
			Accuracy:             2,
			Sprites: []manifest.Sprite{
				{Angle: 0, Width: 32, Height: 32, X: 0},
				{Angle: 45, Width: 32, Height: 32, X: 40},
			},
		},
	}

	sheets := GetSpritesheets(def)
	testSpritesheet(t, &sheets, "32bpp")
	testSpritesheet(t, &sheets, "8bpp")
	testSpritesheet(t, &sheets, "mask")
}

func testSpritesheet(t *testing.T, sheets *Spritesheets, bpp string) {
	sheet, ok := sheets.Data[bpp]

	if !ok {
		t.Fatalf("no " + bpp + "spritesheet present in result")
	}

	expectedRect := image.Rectangle{Max: image.Point{X: 80, Y: 32}}

	if sheet.Image.Bounds() != expectedRect {
		t.Errorf("spritesheet size %v did not match expected size %v", sheet.Image.Bounds(), expectedRect)
	}

	if !imageutils.IsColourEqual(sheet.Image, 79, 0, 65535, 65535, 65535) {
		t.Errorf("blank area of spritesheet not set to white")
	}
}

func Benchmark_32bpp(b *testing.B) {
	spritesheetImage := get32bppSpritesheetImage
	benchmarkSpritesheet(b, spritesheetImage, "32bpp")
}

func Benchmark_8bpp(b *testing.B) {
	spritesheetImage := get8bppSpritesheetImage
	benchmarkSpritesheet(b, spritesheetImage, "8bpp")
}

func benchmarkSpritesheet(b *testing.B, spritesheetImage func(def manifest.Definition, bounds image.Rectangle, spriteInfos []SpriteInfo, depth string) (img image.Image), depth string) {
	object := getObjectForBenchmark("cone.vox", b)
	palette := getPalette(b)

	def := manifest.Definition{
		Object:  object,
		Palette: palette,
		Scale:   2.0,
		Manifest: manifest.Manifest{
			LightingAngle:        45,
			LightingElevation:    60,
			Size:                 object.Size.ToVector3(),
			RenderElevationAngle: 30,
			Accuracy:             2,
			Sprites:              []manifest.Sprite{{Angle: 0, Width: 32, Height: 32, X: 0}},
		},
	}

	bounds := image.Rectangle{Max: image.Point{
		X: int(float64(def.Manifest.Sprites[0].Width) * def.Scale),
		Y: int(float64(def.Manifest.Sprites[0].Height) * def.Scale),
	}}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rect := getSpriteSizeForAngle(def.Manifest.Sprites[0], def.Scale)

		smp := sampler.Disc(rect.Max.X, rect.Max.Y, def.Manifest.Accuracy, 0, 0)
		spr := manifest.Sprite{OffsetX: 0, OffsetY: 0}
		ro := raycaster.GetRaycastOutput(def.Object, def.Manifest, def.Manifest.Sprites[0], smp)
		so := sprite.GetShaderOutput(ro, spr, def, rect.Max.X, rect.Max.Y)

		info := SpriteInfo{
			ShaderOutput: so,
			SpriteBounds: rect,
		}

		spritesheetImage(def, bounds, []SpriteInfo{info}, depth)
	}
}

func getPalette(b *testing.B) colour.Palette {
	pFile, err := os.Open("../../files/ttd_palette.json")
	if err != nil {
		b.Fatalf("Could nopt open palette file: %v", err)
	}

	palette, err := colour.FromJson(pFile)
	if err != nil {
		b.Fatalf("Could not open palette file: %v", err)
	}

	pFile.Close()
	return palette
}

func getObjectForBenchmark(filename string, b *testing.B) voxelobject.ProcessedVoxelObject {
	mv, err := magica.FromFile("../raycaster/testdata/"+filename)
	if err != nil {
		b.Fatalf("error loading test file: %v", err)
	}

	v := voxelobject.GetProcessedVoxelObject(mv, &colour.Palette{}, false, false)
	return v
}
