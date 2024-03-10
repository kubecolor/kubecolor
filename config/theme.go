package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewBaseTheme returns the base color schema depending on the dark/light setting
func NewBaseTheme(preset Preset) *Theme {
	switch preset {
	case PresetDark:
		return &Theme{
			Default: MustParseColor("yellow"),
			Base: ThemeBase{
				Info:      MustParseColor("white"),
				Primary:   MustParseColor("magenta"),
				Secondary: MustParseColor("cyan"),
				Success:   MustParseColor("green"),
				Warning:   MustParseColor("yellow"),
				Danger:    MustParseColor("red"),
			},
		}

	case PresetLight:
		return &Theme{
			Default: MustParseColor("yellow"),
			Base: ThemeBase{
				Info:      MustParseColor("black"),
				Primary:   MustParseColor("magenta"),
				Secondary: MustParseColor("blue"),
				Success:   MustParseColor("green"),
				Warning:   MustParseColor("yellow"),
				Danger:    MustParseColor("red"),
			},
		}

	default:
		panic(fmt.Sprintf("invalid theme preset: %s", preset))
	}
}

// Theme is the root theme config.
type Theme struct {
	// TODO: Rename to more specific
	Default Color // default when no specific mapping is found for the command
	// TODO: Remove in favor of sub-command specific
	Error Color // used when the value is required or is an error

	// TODO: Move to Theme.Table
	DurationFresh Color // color used when the time value is under a certain delay

	Header  Color      // used to print headers
	Columns ColorSlice // used to display multiple colons, cycle between colors

	Base   ThemeBase   // base colors for themes
	Data   ThemeData   // colors for representing data
	Status ThemeStatus // generic status coloring (e.g "Ready", "Terminating")
	Apply  ThemeApply  // used in "kubectl apply"
}

func (t Theme) ApplyViperDefaults(v *viper.Viper) {
	// Base colors are applied first
	t.Base.ApplyViperDefaults(v)

	viperSetDefaultColor(v, "theme.default", t.Default)
	viperSetDefaultColorOrKey(v, "theme.error", t.Error, baseDanger)
	viperSetDefaultColorOrKey(v, "theme.durationfresh", t.DurationFresh, baseSuccess)
	viperSetDefaultColorOrKey(v, "theme.header", t.Header, baseInfo)
	viperSetDefaultColorSliceOrKeys(v, "theme.columns", t.Columns, baseInfo, baseSecondary)

	t.Data.ApplyViperDefaults(v)
	t.Status.ApplyViperDefaults(v)
	t.Apply.ApplyViperDefaults(v)
}

// ThemeBase contains base colors that other theme fields can default to,
// just to make overriding themes easier.
//
// These fields should never be referenced in the printers.
// Instead, they should use the more specific fields, such as [ThemeApply.Created]
type ThemeBase struct {
	Info      Color // general color for when things are informational
	Primary   Color // general color for when things are focus
	Secondary Color // general color for when things are secondary focus
	Success   Color // general color for when things are good
	Warning   Color // general color for when things are wrong
	Danger    Color // general color for when things are bad
}

func (t ThemeBase) ApplyViperDefaults(v *viper.Viper) {
	viperSetDefaultColor(v, baseInfo, t.Info)
	viperSetDefaultColor(v, basePrimary, t.Primary)
	viperSetDefaultColor(v, baseSecondary, t.Secondary)
	viperSetDefaultColor(v, baseSuccess, t.Success)
	viperSetDefaultColor(v, baseWarning, t.Warning)
	viperSetDefaultColor(v, baseDanger, t.Danger)
}

// baseKey are utility strings for referencing the viper keys for the
// base theme colors.
const (
	baseInfo      = "theme.base.info"
	basePrimary   = "theme.base.primary"
	baseSecondary = "theme.base.secondary"
	baseSuccess   = "theme.base.success"
	baseWarning   = "theme.base.warning"
	baseDanger    = "theme.base.danger"
)

// ThemeData holds colors for when representing parsed data.
// Such as in YAML, JSON, and even some "kubectl describe" values
type ThemeData struct {
	String Color // default color for strings
	True   Color // used when value is true
	False  Color // used when value is false
	Number Color // used when the value is a number
	Null   Color // used when the value is null, nil, or none
}

func (t ThemeData) ApplyViperDefaults(v *viper.Viper) {
	viperSetDefaultColorOrKey(v, "theme.data.string", t.String, baseInfo)
	viperSetDefaultColorOrKey(v, "theme.data.true", t.True, baseSuccess)
	viperSetDefaultColorOrKey(v, "theme.data.false", t.False, baseDanger)
	viperSetDefaultColorOrKey(v, "theme.data.number", t.Number, basePrimary)
	viperSetDefaultColorOrKey(v, "theme.data.null", t.Null, baseWarning)
}

// ThemeApply holds colors for the "kubectl apply" output.
type ThemeApply struct {
	Created    Color
	Configured Color
	Unchanged  Color
	DryRun     Color
	Fallback   Color
}

func (t ThemeApply) ApplyViperDefaults(v *viper.Viper) {
	viperSetDefaultColorOrKey(v, "theme.apply.created", t.Created, baseSuccess)
	viperSetDefaultColorOrKey(v, "theme.apply.configured", t.Configured, baseWarning)
	viperSetDefaultColorOrKey(v, "theme.apply.unchanged", t.Unchanged, basePrimary)
	viperSetDefaultColorOrKey(v, "theme.apply.dryrun", t.DryRun, baseSecondary)
	viperSetDefaultColorOrKey(v, "theme.apply.fallback", t.Fallback, baseSuccess)
}

// ThemeStatus holds colors for status texts, used in for example
// the "kubectl get" status column
type ThemeStatus struct {
	Success Color // e.g "Running", "Ready"
	Warning Color // e.g "Terminating"
	Error   Color // e.g "Failed", "Unhealthy"
}

func (t ThemeStatus) ApplyViperDefaults(v *viper.Viper) {
	viperSetDefaultColorOrKey(v, "theme.status.success", t.Success, baseSuccess)
	viperSetDefaultColorOrKey(v, "theme.status.warning", t.Warning, baseWarning)
	viperSetDefaultColorOrKey(v, "theme.status.error", t.Error, baseDanger)
}

func viperSetDefaultColorOrKey(v *viper.Viper, key string, value Color, otherKey string) {
	if viperSetDefaultColor(v, key, value) {
		return
	}
	// We can read from [viper.Viper.Get] here, as it contains the user's
	// configs too.
	// But we cannot just pass in the [ThemeBase] struct, as we have not
	// called [viper.Viper.Unmarshal] yet.
	v.SetDefault(key, v.Get(otherKey))
}

func viperSetDefaultColor(v *viper.Viper, key string, value Color) bool {
	if value == (Color{}) {
		return false
	}
	v.SetDefault(key, value)
	return true
}

func viperSetDefaultColorSliceOrKeys(v *viper.Viper, key string, value []Color, otherKeys ...string) {
	if viperSetDefaultColorSlice(v, key, value) {
		return
	}

	values := make([]any, 0, len(otherKeys))
	for _, k := range otherKeys {
		val := v.Get(k)
		if val != nil {
			values = append(values, val)
		}
	}
	if len(values) > 0 {
		v.SetDefault(key, values)
	}
}

func viperSetDefaultColorSlice(v *viper.Viper, key string, value []Color) bool {
	if len(value) == 0 {
		return false
	}
	v.SetDefault(key, value)
	return true
}
