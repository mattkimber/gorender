package colour

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
)

type PaletteEntry struct {
	R, G, B byte
	Range   *PaletteRange
}

type PaletteRange struct {
	Start                    byte `json:"start"`
	End                      byte `json:"end"`
	IsPrimaryCompanyColour   bool `json:"is_primary_company_colour"`
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

		if entry.Range != nil {
			if entry.Range.IsPrimaryCompanyColour || entry.Range.IsSecondaryCompanyColour {
				y := (19595*uint32(entry.R) + 38470*uint32(entry.G) + 7471*uint32(entry.B) + 1<<15) >> 8
				return y, y, y
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

func GetPaletteFromJson(handle io.Reader) (p Palette, err error) {
	data, err := ioutil.ReadAll(handle)

	if err != nil {
		return Palette{}, err
	}

	if err := json.Unmarshal(data, &p); err != nil {
		return Palette{}, err
	}

	if err := p.SetRanges(p.Ranges); err != nil {
		return Palette{}, err
	}

	return
}

func (p *Palette) SetRanges(ranges []PaletteRange) (err error) {
	p.Ranges = ranges

	for i := range p.Entries {
		p.Entries[i].Range = nil
	}

	for i, r := range ranges {
		for j := r.Start; j <= r.End; j++ {
			if p.Entries[j].Range != nil {
				return fmt.Errorf("range %d overlaps colour %d", i, j)
			}
			p.Entries[j].Range = &ranges[i]
		}
	}

	return nil
}

func (p *Palette) GetFromReader(handle io.Reader) (err error) {
	*p, err = GetPaletteFromJson(handle)
	return err
}
