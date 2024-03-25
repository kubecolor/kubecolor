package config

import (
	"encoding"
	"fmt"
	"strings"
)

type Preset byte

const (
	PresetUnknown Preset = iota
	PresetDark
	PresetLight

	PresetPre0021Dark
	PresetPre0021Light

	PresetProtDark
	PresetProtLight

	PresetDeutDark
	PresetDeutLight

	PresetTritDark
	PresetTritLight
)

var (
	PresetDefault = PresetDark

	// AllPresets is used in places like the internal/cmd/configschema package
	// to show all available options.
	AllPresets = []Preset{
		PresetDark,
		PresetLight,

		PresetPre0021Dark,
		PresetPre0021Light,
	}

	_ encoding.TextMarshaler   = &PresetDefault
	_ encoding.TextUnmarshaler = &PresetDefault
)

func (p Preset) String() string {
	switch p {
	case PresetUnknown:
		return "unknown"
	case PresetDark:
		return "dark"
	case PresetLight:
		return "light"
	case PresetPre0021Dark:
		return "pre-0.0.21-dark"
	case PresetPre0021Light:
		return "pre-0.0.21-light"
	case PresetProtDark:
		return "protanopia-dark"
	case PresetProtLight:
		return "protanopia-light"
	case PresetDeutDark:
		return "deuteranopia-dark"
	case PresetDeutLight:
		return "deuteranopia-light"
	case PresetTritDark:
		return "tritanopia-dark"
	case PresetTritLight:
		return "tritanopia-light"
	default:
		return fmt.Sprintf("%[1]T(%[1]d)", p)
	}
}

func ParsePreset(s string) (Preset, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	// Don't try to parse [PresetUnknown]. It's for internal usage only
	case "dark":
		return PresetDark, nil
	case "light":
		return PresetLight, nil
	case "pre-0.0.21-dark":
		return PresetPre0021Dark, nil
	case "pre-0.0.21-light":
		return PresetPre0021Light, nil
	case "protanopia-dark":
		return PresetProtDark, nil
	case "protanopia-light":
		return PresetProtLight, nil
	case "deuteranopia-dark":
		return PresetDeutDark, nil
	case "deuteranopia-light":
		return PresetDeutLight, nil
	case "tritanopia-dark":
		return PresetTritDark, nil
	case "tritanopia-light":
		return PresetTritLight, nil
	default:
		return Preset(0), fmt.Errorf("invalid theme preset: %q", s)
	}
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
