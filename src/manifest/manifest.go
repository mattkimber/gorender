package manifest

import (
	"encoding/json"
	"geometry"
	"io"
	"io/ioutil"
)

type Sprite struct {
	Angle  float64 `json:"angle"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	X      int
}

type Manifest struct {
	LightingAngle        int            `json:"lighting_angle"`
	LightingElevation    int            `json:"lighting_elevation"`
	Size                 geometry.Point `json:"size"`
	RenderElevationAngle int            `json:"render_elevation"`
	Sprites              []Sprite       `json:"sprites"`
	DepthInfluence       float64        `json:"depth_influence"`
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
