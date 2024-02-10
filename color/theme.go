package color

import (
	"fmt"
	"os"
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
	if err := t.Status.OverrideFromEnv(); err != nil {
		return err
	}
	if err := t.Apply.OverrideFromEnv(); err != nil {
		return err
	}

	if err := setColorFromEnv(&t.DefaultColor, "KUBECOLOR_THEME_DEFAULT"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.ErrorColor, "KUBECOLOR_THEME_ERROR"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.StringColor, "KUBECOLOR_THEME_STRING"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.TrueColor, "KUBECOLOR_THEME_TRUE"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.FalseColor, "KUBECOLOR_THEME_FALSE"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.NumberColor, "KUBECOLOR_THEME_NUMBER"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.NullColor, "KUBECOLOR_THEME_NULL"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.DurationFreshColor, "KUBECOLOR_THEME_DURATION_FRESH"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.HeaderColor, "KUBECOLOR_THEME_HEADER"); err != nil {
		return err
	}
	if err := setColorSliceFromEnv(&t.ColumnColorCycle, "KUBECOLOR_THEME_COLUMN_CYCLE"); err != nil {
		return err
	}
	return nil
}

func (t *ThemeStatus) OverrideFromEnv() error {
	if err := setColorFromEnv(&t.SuccessColor, "KUBECOLOR_THEME_STATUS_SUCCESS"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.WarningColor, "KUBECOLOR_THEME_STATUS_WARNING"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.ErrorColor, "KUBECOLOR_THEME_STATUS_ERROR"); err != nil {
		return err
	}
	return nil
}

func (t *ThemeApply) OverrideFromEnv() error {
	if err := setColorFromEnv(&t.CreatedColor, "KUBECOLOR_THEME_APPLY_CREATED"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.ConfiguredColor, "KUBECOLOR_THEME_APPLY_CONFIGURED"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.UnchangedColor, "KUBECOLOR_THEME_APPLY_UNCHANGED"); err != nil {
		return err
	}
	if err := setColorFromEnv(&t.DryRunColor, "KUBECOLOR_THEME_APPLY_DRYRUN"); err != nil {
		return err
	}
	return nil
}

func setColorSliceFromEnv(target *[]Color, env string) error {
	value := os.Getenv(env)
	if value == "" {
		return nil
	}
	var cols []Color
	// The split character is picked to specifically not collide
	// with gookit/color's syntax: [https://pkg.go.dev/github.com/gookit/color#ParseCodeFromAttr]
	for _, v := range strings.Split(value, "/") {
		col, err := Parse(v)
		if err != nil {
			return fmt.Errorf("%s: %w", env, err)
		}
		cols = append(cols, col)
	}
	*target = cols
	return nil
}

func setColorFromEnv(target *Color, env string) error {
	value := os.Getenv(env)
	if value == "" {
		return nil
	}
	col, err := Parse(value)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}
	*target = col
	return nil
}
