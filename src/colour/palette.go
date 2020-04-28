package colour

import (
	"encoding/json"
	"image/color"
	"io"
	"io/ioutil"
)

type PaletteEntry struct {
	R, G, B byte
}

type PaletteRange struct {
	Start byte `json:"start"`
	End   byte `json:"end"`
	IsPrimaryCompanyColour bool `json:"is_primary_company_colour"`
	IsSecondaryCompanyColour bool `json:"is_secondary_company_colour"`
}

type Palette struct {
	Entries []PaletteEntry `json:"entries"`
	Ranges  []PaletteRange `json:"ranges"`
}

func (p Palette) GetRGB(index byte) (r, g, b uint32) {
	if int(index) < len(p.Entries) {
		entry := p.Entries[index]
		rgba := color.RGBA{R: entry.R, B: entry.B, G: entry.G}
		r, g, b, _ = rgba.RGBA()

		for _, rng := range p.Ranges {
			if index >= rng.Start && index <= rng.End {
				if rng.IsPrimaryCompanyColour || rng.IsSecondaryCompanyColour {
					y := (19595*uint32(entry.R) + 38470*uint32(entry.G) + 7471*uint32(entry.B) + 1<<15) >> 8
					return y,y,y
				}
			}
		}

		return
	}

	return 0, 0, 0
}

func (p *PaletteEntry) UnmarshalJSON(data []byte) error {
	i := []interface{}{&p.R, &p.G, &p.B}

	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	return nil
}

func GetPaletteFromJson(handle io.Reader) (p Palette) {
	data, err := ioutil.ReadAll(handle)

	if err != nil {
		return Palette{}
	}

	if err := json.Unmarshal(data, &p); err != nil {
		return Palette{}
	}

	return
}
