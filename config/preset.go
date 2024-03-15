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
)

var (
	PresetDefault = PresetDark

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
