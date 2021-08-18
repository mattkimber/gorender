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
	IsAnimatedLight          bool `json:"is_animated_light"`
	IsProcessColour          bool `json:"is_process_colour"`
	Smoothness               int  `json:"smoothness"`
	IsNonRenderable          bool `json:"non_renderable"`
	MaxGapInRegion           int  `json:"max_gap_in_region"`
	ExpectedColourRange      byte `json:"expected_colour_range"`
}

type Palette struct {
	Entries                           []PaletteEntry `json:"entries"`
	Ranges                            []PaletteRange `json:"ranges"`
	CompanyColourLightingContribution float64        `json:"company_colour_lighting_contribution"`
	DefaultBrightness                 float64        `json:"default_brightness"`
	CompanyColourLightingScale        float64        `json:"company_colour_lighting_scale"`
}

func (pe *PaletteEntry) GetRGB() (output RGB) {
	output.R = float64(pe.R) * 255
	output.G = float64(pe.G) * 255
	output.B = float64(pe.B) * 255
	return
}

// Get a Go palette
func (p Palette) GetGoPalette() (pal color.Palette) {
	pal = make([]color.Color, len(p.Entries))

	for i, e := range p.Entries {
		pal[i] = color.RGBA{R: e.R, G: e.G, B: e.B, A: 255}
	}

	return
}

// Get the palette of non-special colours
func (p Palette) GetRegularPalette() (pal []RGB) {
	pal = make([]RGB, len(p.Entries))

	for i, e := range p.Entries {
		if !p.IsSpecialColour(byte(i)) {
			pal[i] = FromPaletteEntry(e)
		} else {
			pal[i] = RGB{R: 65535, G: 0, B: 65535}
		}
	}

	return
}

// Get the palette of primary company colours
func (p Palette) GetPrimaryCompanyColourPalette() (pal []RGB) {
	pal = make([]RGB, len(p.Entries))

	for i, e := range p.Entries {
		if e.Range != nil && i != 0 && i != 255 && e.Range.IsPrimaryCompanyColour {
			pal[i] = FromPaletteEntry(e)
		} else {
			pal[i] = RGB{R: 65535, G: 0, B: 65535}
		}
	}

	return
}

// Get the palette of secondary company colours
func (p Palette) GetSecondaryCompanyColourPalette() (pal []RGB) {
	pal = make([]RGB, len(p.Entries))

	for i, e := range p.Entries {
		if e.Range != nil && i != 0 && i != 255 && e.Range.IsSecondaryCompanyColour {
			pal[i] = FromPaletteEntry(e)
		} else {
			pal[i] = RGB{R: 65535, G: 0, B: 65535}
		}
	}

	return
}

// Get the palette of animated colours
func (p Palette) GetAnimatedPalette() (pal []RGB) {
	pal = make([]RGB, len(p.Entries))

	for i, e := range p.Entries {
		if e.Range != nil && i != 0 && i != 255 && e.Range.IsAnimatedLight {
			pal[i] = FromPaletteEntry(e)
		} else {
			pal[i] = RGB{R: 65535, G: 0, B: 65535}
		}
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
			if entry.Range.IsPrimaryCompanyColour || entry.Range.IsSecondaryCompanyColour || entry.Range.IsAnimatedLight {
				return index
			}
		}
	}

	return
}

func (p Palette) IsRenderable(index byte) bool {
	if int(index) < len(p.Entries) && p.Entries[index].Range != nil {
		return !p.Entries[index].Range.IsNonRenderable
	}

	return false
}

func (p Palette) IsSpecialColour(index byte) bool {
	if int(index) < len(p.Entries) && p.Entries[index].Range != nil {
		return p.Entries[index].Range.IsPrimaryCompanyColour || p.Entries[index].Range.IsSecondaryCompanyColour || p.Entries[index].Range.IsAnimatedLight || p.Entries[index].Range.IsNonRenderable
	}

	return false
}

func (p Palette) GetRGB(index byte, resolveSpecialColours bool) (output RGB) {
	if int(index) < len(p.Entries) {
		entry := p.Entries[index]
		rgba := color.RGBA{R: entry.R, B: entry.B, G: entry.G}
		r32, g32, b32, _ := rgba.RGBA()
		output = RGB{
			R: float64(r32),
			G: float64(g32),
			B: float64(b32),
		}

		if !resolveSpecialColours {
			return
		}

		if entry.Range != nil {
			if entry.Range.IsPrimaryCompanyColour || entry.Range.IsSecondaryCompanyColour {
				cc := float64((19595*uint32(entry.R) + 38470*uint32(entry.G) + 7471*uint32(entry.B) + 1<<15) >> 8)
				y := (p.DefaultBrightness * 32767.0 * (1 - p.CompanyColourLightingContribution)) + (cc * p.CompanyColourLightingContribution)
				return RGB{R: y, G: y, B: y}
			}

			if entry.Range.IsAnimatedLight {
				return RGB{R: 22000, G: 22000, B: 22000}
			}
		}

		return
	}

	return RGB{R: 0, G: 0, B: 0}
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

func (p Palette) GetLitRGB(index byte, l float64, brightness float64, contrast float64, resolveSpecialColours bool, influence float64) (output RGB) {
	output = p.GetRGB(index, resolveSpecialColours)

	entry := p.Entries[index]
	if resolveSpecialColours && entry.Range != nil && (entry.Range.IsPrimaryCompanyColour || entry.Range.IsSecondaryCompanyColour) {
		l = l * p.CompanyColourLightingScale
	}

	if entry.Range != nil && entry.Range.IsAnimatedLight {
		l = 0.5
	}

	// Clamp to [-1,1]
	if l > 1 {
		l = 1
	} else if l < -1 {
		l = -1
	}

	if l >= 0 {
		// interpolate towards white
		output.R = (output.R * (1 - l)) + (65535 * l)
		output.G = (output.G * (1 - l)) + (65535 * l)
		output.B = (output.B * (1 - l)) + (65535 * l)
	} else if l < 0 {
		// interpolate towards black
		output.R = output.R * (1 + l)
		output.G = output.G * (1 + l)
		output.B = output.B * (1 + l)
	}

	// Apply brightness/contrast
	output.R += brightness
	output.G += brightness
	output.B += brightness

	output.R = (contrast*(output.R-32767) + 32767) * influence
	output.G = (contrast*(output.G-32767) + 32767) * influence
	output.B = (contrast*(output.B-32767) + 32767) * influence

	return
}

func (pe *PaletteEntry) UnmarshalJSON(data []byte) error {
	i := []interface{}{&pe.R, &pe.G, &pe.B}

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

		// Set the default for max region gap
		if r.MaxGapInRegion == 0 {
			ranges[i].MaxGapInRegion = 6
			ranges[i].ExpectedColourRange = 3
		}

		for j := int(r.Start); j <= int(r.End); j++ {
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
