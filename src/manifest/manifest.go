package manifest

import (
	"colour"
	"encoding/json"
	"geometry"
	"io"
	"io/ioutil"
	"voxelobject"
)

type Definition struct {
	Object   voxelobject.ProcessedVoxelObject
	Palette  colour.Palette
	Manifest Manifest
	Scale    float64
	Debug    bool
	Time     bool
}

type Sprite struct {
	Angle  float64 `json:"angle"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	X      int
	Flip   bool `json:"flip"`
}

type Manifest struct {
	LightingAngle        int              `json:"lighting_angle"`
	LightingElevation    int              `json:"lighting_elevation"`
	Size                 geometry.Vector3 `json:"size"`
	RenderElevationAngle int              `json:"render_elevation"`
	Sprites              []Sprite         `json:"sprites"`
	DepthInfluence       float64          `json:"depth_influence"`
	TiledNormals         bool             `json:"tiled_normals"`
	SoftenEdges          float64          `json:"soften_edges"`
}

func FromJson(handle io.Reader) (manifest Manifest, err error) {
	data, err := ioutil.ReadAll(handle)

	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &manifest); err != nil {
		return
	}

	return
}

func (m *Manifest) GetFromReader(handle io.Reader) (err error) {
	*m, err = FromJson(handle)
	return err
}

func (d *Definition) SoftenEdges() bool {
	return d.Scale >= d.Manifest.SoftenEdges
}
