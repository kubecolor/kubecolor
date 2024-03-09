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
