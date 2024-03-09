package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Theme struct {
	Default Color // default when no specific mapping is found for the command
	Error   Color // used when the value is required or is an error

	String        Color // default color for strings
	True          Color // used when value is true
	False         Color // used when value is false
	Number        Color // used when the value is a number
	Null          Color // used when the value is null, nil, or none
	DurationFresh Color // color used when the time value is under a certain delay

	Header      Color   // used to print headers
	ColumnCycle []Color // used to display multiple colons, cycle between colors

	Status ThemeStatus // generic status coloring (e.g "Ready", "Terminating")
	Apply  ThemeApply  // used in "kubectl apply"
}

func (t Theme) ApplyViperDefaults(v *viper.Viper) {
	v.SetDefault("theme.default", t.Default)
	v.SetDefault("theme.error", t.Error)
	v.SetDefault("theme.string", t.String)
	v.SetDefault("theme.true", t.True)
	v.SetDefault("theme.false", t.False)
	v.SetDefault("theme.number", t.Number)
	v.SetDefault("theme.null", t.Null)
	v.SetDefault("theme.durationfresh", t.DurationFresh)
	v.SetDefault("theme.header", t.Header)
	v.SetDefault("theme.columncycle", t.ColumnCycle)

	t.Status.ApplyViperDefaults(v)
	t.Apply.ApplyViperDefaults(v)
}

type ThemeApply struct {
	Created    Color
	Configured Color
	Unchanged  Color
	DryRun     Color
	Fallback   Color
}

func (t ThemeApply) ApplyViperDefaults(v *viper.Viper) {
	v.SetDefault("theme.apply.created", t.Created)
	v.SetDefault("theme.apply.configured", t.Configured)
	v.SetDefault("theme.apply.unchanged", t.Unchanged)
	v.SetDefault("theme.apply.dryrun", t.DryRun)
}

type ThemeStatus struct {
	Success Color // e.g "Running", "Ready"
	Warning Color // e.g "Terminating"
	Error   Color // e.g "Failed", "Unhealthy"
}

func (t ThemeStatus) ApplyViperDefaults(v *viper.Viper) {
	v.SetDefault("theme.status.success", t.Success)
	v.SetDefault("theme.status.warning", t.Warning)
	v.SetDefault("theme.status.error", t.Error)
}

// NewTheme returns the base color schema depending on the dark/light setting
func NewTheme(preset Preset) *Theme {
	switch preset {
	case PresetDark:
		return &Theme{
			Default:       MustParseColor("yellow"),
			String:        MustParseColor("white"),
			True:          MustParseColor("green"),
			False:         MustParseColor("red"),
			Number:        MustParseColor("magenta"),
			Null:          MustParseColor("yellow"),
			Header:        MustParseColor("white"),
			DurationFresh: MustParseColor("green"),
			Error:         MustParseColor("red"),
			ColumnCycle:   []Color{MustParseColor("white"), MustParseColor("cyan")},
			Status: ThemeStatus{
				Success: MustParseColor("green"),
				Warning: MustParseColor("yellow"),
				Error:   MustParseColor("red"),
			},
			Apply: ThemeApply{
				Created:    MustParseColor("green"),
				Configured: MustParseColor("yellow"),
				Unchanged:  MustParseColor("magenta"),
				DryRun:     MustParseColor("cyan"),
				Fallback:   MustParseColor("green"),
			},
		}

	case PresetLight:
		return &Theme{
			Default:       MustParseColor("yellow"),
			String:        MustParseColor("black"),
			True:          MustParseColor("green"),
			False:         MustParseColor("red"),
			Number:        MustParseColor("magenta"),
			Null:          MustParseColor("yellow"),
			Header:        MustParseColor("black"),
			DurationFresh: MustParseColor("green"),
			Error:         MustParseColor("red"),
			ColumnCycle:   []Color{MustParseColor("black"), MustParseColor("blue")},
			Status: ThemeStatus{
				Success: MustParseColor("green"),
				Warning: MustParseColor("yellow"),
				Error:   MustParseColor("red"),
			},
			Apply: ThemeApply{
				Created:    MustParseColor("green"),
				Configured: MustParseColor("yellow"),
				Unchanged:  MustParseColor("magenta"),
				DryRun:     MustParseColor("blue"),
				Fallback:   MustParseColor("green"),
			},
		}

	default:
		panic(fmt.Sprintf("invalid theme preset: %s", preset))
	}
}
