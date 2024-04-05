package config

import (
	"encoding"
	"fmt"
	"strings"
)

type Preset string

const (
	// NOTE: When adding presets, remember to add them to [AllPresets] slice too.

	// "Zero value", i.e empty theme selected
	PresetNone Preset = ""

	// Default themes
	PresetDark  Preset = "dark"
	PresetLight Preset = "light"

	// Color blind focused themes
	PresetProtDark  Preset = "protanopia-dark"
	PresetProtLight Preset = "protanopia-light"
	PresetDeutDark  Preset = "deuteranopia-dark"
	PresetDeutLight Preset = "deuteranopia-light"
	PresetTritDark  Preset = "tritanopia-dark"
	PresetTritLight Preset = "tritanopia-light"

	// Pre-v0.3.0
	PresetPre030Dark  Preset = "pre-0.3.0-dark"
	PresetPre030Light Preset = "pre-0.3.0-light"

	// Pre-v0.0.21
	PresetPre0021Dark  Preset = "pre-0.0.21-dark"
	PresetPre0021Light Preset = "pre-0.0.21-light"
)

var (
	PresetDefault = PresetDark

	// AllPresets is used in parsing and places like the
	// internal/cmd/configschema package to show all available options.
	AllPresets = []Preset{
		// "Zero value", i.e empty theme selected
		PresetNone,

		// Default themes
		PresetDark,
		PresetLight,

		// Color blind focused themes
		PresetProtDark,
		PresetProtLight,
		PresetDeutDark,
		PresetDeutLight,
		PresetTritDark,
		PresetTritLight,

		// Pre-v0.3.0
		PresetPre030Dark,
		PresetPre030Light,
		// Pre-v0.0.21
		PresetPre0021Dark,
		PresetPre0021Light,
	}

	_ encoding.TextMarshaler   = PresetDefault
	_ encoding.TextUnmarshaler = &PresetDefault
)

func (p Preset) String() string {
	if p == "" {
		return "none"
	}
	return string(p)
}

func ParsePreset(s string) (Preset, error) {
	if s == "" {
		return PresetNone, nil
	}
	maybeValidPreset := Preset(strings.ToLower(s))
	for _, p := range AllPresets {
		if maybeValidPreset == p {
			return p, nil // reuse the interned string
		}
	}
	return PresetNone, fmt.Errorf("invalid theme preset: %q", s)
}

// MarshalText implements [encoding.TextMarshaler].
func (p Preset) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (p *Preset) UnmarshalText(text []byte) error {
	newPreset, err := ParsePreset(string(text))
	if err != nil {
		return err
	}
	*p = newPreset
	return nil
}
