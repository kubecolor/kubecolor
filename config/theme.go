package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kubecolor/kubecolor/internal/stringutil"
	"github.com/spf13/viper"
)

// NewBaseTheme returns the base color schema depending on the dark/light setting
func NewBaseTheme(preset Preset) *Theme {
	switch preset {
	case PresetDark:
		return &Theme{
			Default: MustParseColor("yellow"),
			Base: ThemeBase{
				Key:       MustParseColorSlice("yellow / white"),
				Info:      MustParseColor("white"),
				Primary:   MustParseColor("magenta"),
				Secondary: MustParseColor("cyan"),
				Success:   MustParseColor("green"),
				Warning:   MustParseColor("yellow"),
				Danger:    MustParseColor("red"),
				// TODO: Only set this to "yellow" in "legacy" theme, and instead use "gray:italic" here
				Muted: MustParseColor("yellow"),
			},
			Options: ThemeOptions{
				// Overriding here to make tests happy
				// TODO: Only set this to "yellow" in "legacy" theme, and instead rely on defaultFrom
				Flag: MustParseColor("yellow"),
			},
		}

	case PresetLight:
		return &Theme{
			Default: MustParseColor("yellow"),
			Base: ThemeBase{
				Key:       MustParseColorSlice("yellow / black"),
				Info:      MustParseColor("black"),
				Primary:   MustParseColor("magenta"),
				Secondary: MustParseColor("blue"),
				Success:   MustParseColor("green"),
				Warning:   MustParseColor("yellow"),
				Danger:    MustParseColor("red"),
				// TODO: Only set this to "yellow" in "legacy" theme, and instead use "gray:italic" here
				Muted: MustParseColor("yellow"),
			},
			Options: ThemeOptions{
				// Overriding here to make tests happy
				// TODO: Only set this to "yellow" in "legacy" theme, and instead rely on defaultFrom
				Flag: MustParseColor("yellow"),
			},
		}

	default:
		panic(fmt.Sprintf("invalid theme preset: %s", preset))
	}
}

// Theme is the root theme config.
type Theme struct {
	// Base colors must be first so they're applied first
	Base ThemeBase // base colors for themes

	Default Color // default when no specific mapping is found for the command

	Data   ThemeData   // colors for representing data
	Status ThemeStatus // generic status coloring (e.g "Ready", "Terminating")
	Table  ThemeTable  // used in table output, e.g "kubectl get" and parts of "kubectl describe"
	Stderr ThemeStderr // used in kubectl's stderr output

	Describe ThemeDescribe // used in "kubectl describe"
	Apply    ThemeApply    // used in "kubectl apply"
	Explain  ThemeExplain  // used in "kubectl explain"
	Options  ThemeOptions  // used in "kubectl options"
	Version  ThemeVersion  // used in "kubectl version"
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
	Muted     Color // general color for when things are less relevant

	Key ColorSlice `defaultFromMany:"theme.base.secondary"` // general color for keys
}

// ThemeData holds colors for when representing parsed data.
// Such as in YAML, JSON, and even some "kubectl describe" values
type ThemeData struct {
	Key    ColorSlice `defaultFrom:"theme.base.key"`     // used for the key
	String Color      `defaultFrom:"theme.base.info"`    // used when value is a string
	True   Color      `defaultFrom:"theme.base.success"` // used when value is true
	False  Color      `defaultFrom:"theme.base.danger"`  // used when value is false
	Number Color      `defaultFrom:"theme.base.primary"` // used when the value is a number
	Null   Color      `defaultFrom:"theme.base.muted"`   // used when the value is null, nil, or none

	Duration      Color ``                                 // used when the value is a duration, e.g "12m"
	DurationFresh Color `defaultFrom:"theme.base.success"` // color used when the time value is under a certain delay

	Ratio ThemeDataRatio
}

type ThemeDataRatio struct {
	Zero    Color `defaultFrom:"theme.base.muted"`   // used for "0/0"
	Equal   Color ``                                 // used for "n/n", e.g "1/1"
	Unequal Color `defaultFrom:"theme.base.warning"` // used for "n/m", e.g "0/1"
}

// ThemeStatus holds colors for status texts, used in for example
// the "kubectl get" status column
type ThemeStatus struct {
	Success Color `defaultFrom:"theme.base.success"` // used in status keywords, e.g "Running", "Ready"
	Warning Color `defaultFrom:"theme.base.warning"` // used in status keywords, e.g "Terminating"
	Error   Color `defaultFrom:"theme.base.danger"`  // used in status keywords, e.g "Failed", "Unhealthy"
}

// ThemeTable holds colors for table output
type ThemeTable struct {
	Header  Color      `defaultFrom:"theme.base.info"`                          // used on table headers
	Columns ColorSlice `defaultFromMany:"theme.base.info,theme.base.secondary"` // used on table columns when no other coloring applies such as status or duration coloring. The multiple colors are cycled based on column ID, from left to right.
}

// ThemeStderr holds generic colors for kubectl's stderr output.
type ThemeStderr struct {
	Default Color `defaultFrom:"theme.base.info"`   // default when no specific mapping is found for the output line
	Error   Color `defaultFrom:"theme.base.danger"` // e.g when text contains "error"
}

// ThemeApply holds colors for the "kubectl apply" output.
type ThemeDescribe struct {
	Key ColorSlice `defaultFrom:"theme.base.key"` // used on keys. The multiple colors are cycled based on indentation.
}

// ThemeApply holds colors for the "kubectl apply" output.
type ThemeApply struct {
	Created    Color `defaultFrom:"theme.base.success"`   // used on "deployment.apps/foo created"
	Configured Color `defaultFrom:"theme.base.warning"`   // used on "deployment.apps/bar configured"
	Unchanged  Color `defaultFrom:"theme.base.primary"`   // used on "deployment.apps/quux unchanged"
	DryRun     Color `defaultFrom:"theme.base.secondary"` // used on "deployment.apps/quux created (dry-run)"
	Fallback   Color `defaultFrom:"theme.base.success"`   // used when "kubectl apply" outputs unknown format
}

// ThemeExplain holds colors for the "kubectl explain" output.
type ThemeExplain struct {
	Key      ColorSlice `defaultFrom:"theme.base.key"`    // used on keys. The multiple colors are cycled based on indentation.
	Required Color      `defaultFrom:"theme.base.danger"` // used on the trailing "-required-" string
}

// ThemeOptions holds colors for the "kubectl options" output.
type ThemeOptions struct {
	Flag Color `defaultFrom:"theme.base.secondary"` // e.g "--kubeconfig"
}

// ThemeVersion holds colors for the "kubectl version" output.
type ThemeVersion struct {
	Key ColorSlice `defaultFrom:"theme.base.key"` // used on the key
}

func applyViperDefaults(theme *Theme, v *viper.Viper) {
	themeVal := reflect.ValueOf(theme).Elem()
	themeViperVisitor{
		viper: v,
	}.walkFields(themeVal, "theme")
}

type themeViperVisitor struct {
	viper *viper.Viper
}

func (t themeViperVisitor) walkFields(val reflect.Value, viperKey string) {
	typ := val.Type()
	for i := range val.NumField() {
		fieldTyp := typ.Field(i)
		if fieldTyp.Anonymous || !fieldTyp.IsExported() {
			continue
		}
		fieldVal := val.Field(i)
		// e.g "theme" + field "Default" => "theme.default"
		newViperKey := fmt.Sprintf("%s.%s", viperKey, strings.ToLower(fieldTyp.Name))
		// Only dig deeper if its a theme struct, e.g ThemeApply
		if strings.HasPrefix(fieldTyp.Type.Name(), "Theme") {
			t.walkFields(fieldVal, newViperKey)
			continue
		}
		fieldValInterface := fieldVal.Interface()
		t.visitField(newViperKey, fieldValInterface, fieldTyp.Tag)
	}
}

func (t themeViperVisitor) visitField(viperKey string, value any, tags reflect.StructTag) {
	switch value := value.(type) {
	case Color:
		if _, ok := tags.Lookup("defaultFromMany"); ok {
			panic(fmt.Errorf("%s: cannot use defaultFromMany tag on a Color field", viperKey))
		}
		if defaultFrom, ok := tags.Lookup("defaultFrom"); ok {
			t.setColorOrKey(viperKey, value, defaultFrom)
		} else {
			t.setColor(viperKey, value)
		}
	case ColorSlice:
		if defaultFrom, ok := tags.Lookup("defaultFrom"); ok {
			t.setColorSliceOrKey(viperKey, value, defaultFrom)
		} else if defaultFromMany, ok := tags.Lookup("defaultFromMany"); ok {
			split := stringutil.SplitAndTrimSpace(defaultFromMany, ",")
			t.setColorSliceOrManyKeys(viperKey, value, split)
		} else {
			t.setColorSlice(viperKey, value)
		}
	default:
		panic(fmt.Errorf("%s: unsupported field type: %T", viperKey, value))
	}
}

func (t themeViperVisitor) setColorOrKey(key string, value Color, otherKey string) {
	if t.setColor(key, value) {
		return
	}
	t.viper.SetDefault(key, t.viper.Get(otherKey))
}

func (t themeViperVisitor) setColor(key string, value Color) bool {
	if value == (Color{}) {
		t.viper.SetDefault(key, Color{})
		return false
	}
	t.viper.SetDefault(key, value)
	return true
}

func (t themeViperVisitor) setColorSliceOrKey(key string, value ColorSlice, otherKey string) {
	if t.setColorSlice(key, value) {
		return
	}
	t.viper.SetDefault(key, t.viper.Get(otherKey))
}

func (t themeViperVisitor) setColorSliceOrManyKeys(key string, value ColorSlice, otherKeys []string) {
	if t.setColorSlice(key, value) {
		return
	}
	values := make(ColorSlice, 0, len(otherKeys))
	for _, k := range otherKeys {
		val := t.viper.Get(k)
		col, ok := val.(Color)
		if !ok {
			col = MustParseColor(fmt.Sprint(val))
		}
		if val != nil {
			values = append(values, col)
		}
	}
	if len(values) > 0 {
		t.viper.SetDefault(key, values)
	}
}

func (t themeViperVisitor) setColorSlice(key string, value ColorSlice) bool {
	if len(value) == 0 {
		t.viper.SetDefault(key, ColorSlice{})
		return false
	}
	t.viper.SetDefault(key, value)
	return true
}
