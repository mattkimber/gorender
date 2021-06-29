package manifest

import (
	"encoding/json"
	"github.com/mattkimber/gorender/internal/colour"
	"github.com/mattkimber/gorender/internal/geometry"
	"github.com/mattkimber/gorender/internal/voxelobject"
	"io"
	"io/ioutil"
	"math"
)

type Definition struct {
	Object   voxelobject.ProcessedVoxelObject
	Palette  colour.Palette
	Manifest Manifest
	Scale    float64
	Debug    bool
	Time     bool
	Only8bpp bool
}

type Sprite struct {
	Angle                float64 `json:"angle"`
	Width                int     `json:"width"`
	Height               int     `json:"height"`
	OffsetX              float64 `json:"offset_x"`
	OffsetY              float64 `json:"offset_y"`
	X                    int
	ZError               float64
	Flip                 bool `json:"flip"`
	Slice                int  `json:"slice"`
	RenderElevationAngle int  `json:"render_elevation"`
	Joggle               float64 `json:"joggle"`
}

type Manifest struct {
	LightingAngle        int              `json:"lighting_angle"`
	LightingElevation    int              `json:"lighting_elevation"`
	Size                 geometry.Vector3 `json:"size"`
	RenderElevationAngle int              `json:"render_elevation"`
	Sprites              []Sprite         `json:"sprites"`
	DepthInfluence       float64          `json:"depth_influence"`
	TiledNormals         bool             `json:"tiled_normals"`
	SolidBase            bool             `json:"solid_base"`
	SoftenEdges          float64          `json:"soften_edges"`
	Accuracy             int              `json:"accuracy"`
	Sampler              string           `json:"sampler"`
	Overlap              float64          `json:"overlap"`
	Brightness           float64          `json:"brightness"`
	Contrast             float64          `json:"contrast"`
	DetailBoost          float64          `json:"detail_boost"`
	FadeToBlack          bool             `json:"fade_to_black"`
	EdgeThreshold        float64          `json:"alpha_edge_threshold"`
	HardEdgeThreshold    float64          `json:"hard_edge_threshold"`
	PadToFullLength      bool             `json:"pad_to_full_length"`
	SliceThreshold       int              `json:"slice_threshold"`
	SliceLength          int              `json:"slice_length"`
	SliceOverlap         int              `json:"slice_overlap"`
	Falloff              float64          `json:"falloff_adjustment"`
	RecoveredVoxelSuppression float64     `json:"recovered_voxel_suppression"`
	Joggle               float64 		  `json:"joggle"`
	MaxPush              int              `json:"max_colour_push"`
	Fosterise            bool             `json:"fosterise"`
	NoEdgeFosterisation  bool             `json:"suppress_edge_fosterisation"`
}

func FromJson(handle io.Reader) (manifest Manifest, err error) {
	// Set defaults
	manifest.Accuracy = 2
	manifest.EdgeThreshold = 0.5

	data, err := ioutil.ReadAll(handle)

	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &manifest); err != nil {
		return
	}

	// Convert to standard values
	manifest.Brightness = manifest.Brightness * 65535
	manifest.Contrast += 1.0

	// Set up sprite sizes
	manifest.SetSpriteSizes()

	return
}

func (m *Manifest) GetFromReader(handle io.Reader) (err error) {
	*m, err = FromJson(handle)
	return err
}

func (d *Definition) SoftenEdges() bool {
	return d.Scale >= d.Manifest.SoftenEdges
}

func (m *Manifest) SetSpriteSizes() {
	// Set any auto-height sprites
	for i, _ := range m.Sprites {
		// 0 means "auto"
		if m.Sprites[i].Height == 0 {
			height, delta := getCalculatedSpriteHeight(m, m.Sprites[i])
			m.Sprites[i].Height = height
			m.Sprites[i].ZError = delta
		}

		if m.Sprites[i].RenderElevationAngle == 0 {
			m.Sprites[i].RenderElevationAngle = m.RenderElevationAngle
		}
	}
}

func getCalculatedSpriteHeight(m *Manifest, spr Sprite) (height int, delta float64) {
	size := m.Size
	cos, sin := math.Cos(geometry.DegToRad(spr.Angle)), math.Sin(geometry.DegToRad(spr.Angle))

	xComponent := math.Abs(size.X * cos)
	yComponent := math.Abs(size.Y * sin)

	planeXComponent := math.Abs(size.X * sin)
	planeYComponent := math.Abs(size.Y * cos)

	horizontalSize := (xComponent + yComponent) * math.Sin(geometry.DegToRad(float64(m.RenderElevationAngle)))

	ratio := (horizontalSize + size.Z) / (planeXComponent + planeYComponent)
	spriteSize := ratio * float64(spr.Width)

	spriteSizeRounded := math.Ceil(spriteSize)
	delta = (spriteSizeRounded - spriteSize) / spriteSizeRounded

	return int(spriteSizeRounded), delta
}
