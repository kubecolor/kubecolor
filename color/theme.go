package color

import (
	"fmt"
	"strings"
)

type Theme struct {
	Preset Preset

	DefaultColor Color // default when no specific mapping is found for the command
	ErrorColor   Color // used when the value is required or is an error

	StringColor        Color // default color for strings
	TrueColor          Color // used when value is true
	FalseColor         Color // used when value is false
	NumberColor        Color // used when the value is a number
	NullColor          Color // used when the value is null, nil, or none
	DurationFreshColor Color // color used when the time value is under a certain delay

	HeaderColor      Color   // used to print headers
	ColumnColorCycle []Color // used to display multiple colons, cycle between colors

	Status ThemeStatus // generic status coloring (e.g "Ready", "Terminating")
	Apply  ThemeApply  // used in "kubectl apply"
}

type ThemeApply struct {
	CreatedColor    Color
	ConfiguredColor Color
	UnchangedColor  Color
	DryRunColor     Color
}

type ThemeStatus struct {
	SuccessColor Color // e.g "Running", "Ready"
	WarningColor Color // e.g "Terminating"
	ErrorColor   Color // e.g "Failed", "Unhealthy"
}

type Preset byte

const (
	PresetDefault = PresetDark

	PresetUnknown Preset = iota
	PresetDark
	PresetLight
)

func (t Preset) String() string {
	switch t {
	case PresetUnknown:
		return "unknown"
	case PresetDark:
		return "dark"
	case PresetLight:
		return "light"
	default:
		return fmt.Sprintf("%[1]T(%[1]d)", t)
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

// NewTheme returns the base color schema depending on the dark/light setting
func NewTheme(preset Preset) *Theme {
	switch preset {
	case PresetDark:
		return &Theme{
			Preset:             preset,
			DefaultColor:       Yellow,
			StringColor:        White,
			TrueColor:          Green,
			FalseColor:         Red,
			NumberColor:        Magenta,
			NullColor:          Yellow,
			HeaderColor:        White,
			DurationFreshColor: Green,
			ErrorColor:         Red,
			ColumnColorCycle:   []Color{White, Cyan},
			Status: ThemeStatus{
				SuccessColor: Green,
				WarningColor: Yellow,
				ErrorColor:   Red,
			},
			Apply: ThemeApply{
				CreatedColor:    Green,
				ConfiguredColor: Yellow,
				UnchangedColor:  Magenta,
				DryRunColor:     Cyan,
			},
		}

	case PresetLight:
		return &Theme{
			Preset:             preset,
			DefaultColor:       Yellow,
			StringColor:        Black,
			TrueColor:          Green,
			FalseColor:         Red,
			NumberColor:        Magenta,
			NullColor:          Yellow,
			HeaderColor:        Black,
			DurationFreshColor: Green,
			ErrorColor:         Red,
			ColumnColorCycle:   []Color{Black, Blue},
			Status: ThemeStatus{
				SuccessColor: Green,
				WarningColor: Yellow,
				ErrorColor:   Red,
			},
			Apply: ThemeApply{
				CreatedColor:    Green,
				ConfiguredColor: Yellow,
				UnchangedColor:  Magenta,
				DryRunColor:     Blue,
			},
		}

	default:
		panic(fmt.Sprintf("invalid theme preset: %s", preset))
	}
}

func (t *Theme) OverrideFromEnv() error {
	// TODO: implement this
	return nil
}
