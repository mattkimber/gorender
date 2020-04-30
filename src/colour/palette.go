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
	Smoothness               int  `json:"smoothness"`
}

type Palette struct {
	Entries []PaletteEntry `json:"entries"`
	Ranges  []PaletteRange `json:"ranges"`
}

// Get a Go palette
func (p Palette) GetGoPalette() (pal color.Palette) {
	pal = make([]color.Color, len(p.Entries))

	for i, e := range p.Entries {
		pal[i] = color.RGBA{R: e.R, G: e.G, B: e.B, A: 255}
	}

	return
}

func (p Palette) GetSmoothness(index byte) (smoothness int) {
	if int(index) < len(p.Entries) && p.Entries[index].Range != nil {
		smoothness = p.Entries[index].Range.Smoothness
	}

	return
}

func (p Palette) GetMaskColour(index byte) (msk byte) {
	if int(index) < len(p.Entries) {
		entry := p.Entries[index]
		if entry.Range != nil {
			if entry.Range.IsPrimaryCompanyColour || entry.Range.IsSecondaryCompanyColour {
				return index
			}
		}
	}

	return
}

func (p Palette) IsCompanyColour(index byte) bool {
	if int(index) < len(p.Entries) && p.Entries[index].Range != nil {
		return p.Entries[index].Range.IsPrimaryCompanyColour || p.Entries[index].Range.IsSecondaryCompanyColour
	}

	return false
}

func (p Palette) GetRGB(index byte) (r, g, b uint16) {
	if int(index) < len(p.Entries) {
		entry := p.Entries[index]
		rgba := color.RGBA{R: entry.R, B: entry.B, G: entry.G}
		r32, g32, b32, _ := rgba.RGBA()
		r, g, b = uint16(r32), uint16(g32), uint16(b32)

		if entry.Range != nil {
			if entry.Range.IsPrimaryCompanyColour || entry.Range.IsSecondaryCompanyColour {
				y := uint16((19595*uint32(entry.R) + 38470*uint32(entry.G) + 7471*uint32(entry.B) + 1<<15) >> 8)
				return y, y, y
			}
		}

		return
	}

	return 0, 0, 0
}

func (p Palette) GetLitRGB(index byte, l float64) (r, g, b uint16) {
	r, g, b = p.GetRGB(index)

	// clamp to [-1,1]
	if l > 1 {
		l = 1
	} else if l < -1 {
		l = -1
	}

	if l > 0 {
		// interpolate towards white
		r = uint16((float64(r) * (1 - l)) + (65535 * l))
		g = uint16((float64(g) * (1 - l)) + (65535 * l))
		b = uint16((float64(b) * (1 - l)) + (65535 * l))
	} else if l < 0 {
		// interpolate towards black
		r = uint16(float64(r) * (1 + l))
		g = uint16(float64(g) * (1 + l))
		b = uint16(float64(b) * (1 + l))
	}

	return
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
