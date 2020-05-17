package colour

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"math"
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
	Entries                           []PaletteEntry `json:"entries"`
	Ranges                            []PaletteRange `json:"ranges"`
	CompanyColourLightingContribution float64        `json:"company_colour_lighting_contribution"`
	DefaultBrightness                 float64        `json:"default_brightness"`
	CompanyColourLightingScale        float64        `json:"company_colour_lighting_scale"`
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
				cc := float64((19595*uint32(entry.R) + 38470*uint32(entry.G) + 7471*uint32(entry.B) + 1<<15) >> 8)
				y := uint16((p.DefaultBrightness * 32767.0 * (1 - p.CompanyColourLightingContribution)) + (cc * p.CompanyColourLightingContribution))
				return y, y, y
			}
		}

		return
	}

	return 0, 0, 0
}

func (p Palette) GetLitIndexed(index byte, l float64) (idx byte) {
	if int(index) < len(p.Entries) {
		rng := p.Entries[index].Range
		if rng != nil {
			min, max := rng.Start, rng.End
			spread := max - min
			offsetIndex := float64(index) + math.Round(float64(spread)*(l/2))

			if offsetIndex < float64(min) {
				return min
			}
			if offsetIndex > float64(max) {
				return max
			}

			return byte(offsetIndex)
		}
	}

	return index
}

func (p Palette) GetLitRGB(index byte, l float64) (r, g, b uint16) {
	r, g, b = p.GetRGB(index)

	entry := p.Entries[index]
	if entry.Range != nil && (entry.Range.IsPrimaryCompanyColour || entry.Range.IsSecondaryCompanyColour) {
		l = l * p.CompanyColourLightingScale
	}

	// clamp to [-1,1]
	if l > 1 {
		l = 1
	} else if l < -1 {
		l = -1
	}

	if l >= 0 {
		// interpolate towards white
		r = clamp((float64(r) * (1 - l)) + (65535 * l))
		g = clamp((float64(g) * (1 - l)) + (65535 * l))
		b = clamp((float64(b) * (1 - l)) + (65535 * l))
	} else if l < 0 {
		// interpolate towards black
		r = clamp(float64(r) * (1 + l))
		g = clamp(float64(g) * (1 + l))
		b = clamp(float64(b) * (1 + l))
	}

	return
}

func clamp(input float64) uint16 {
	if input > 65535 {
		return 65535
	} else if input < 256 {
		return 256
	}

	return uint16(input)
}

func (p *PaletteEntry) UnmarshalJSON(data []byte) error {
	i := []interface{}{&p.R, &p.G, &p.B}

	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	return nil
}

func FromJson(handle io.Reader) (p Palette, err error) {
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
	*p, err = FromJson(handle)
	return err
}
