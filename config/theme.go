package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kubecolor/kubecolor/config/color"
	"github.com/kubecolor/kubecolor/internal/stringutil"
	"github.com/spf13/viper"
)

// NewBaseTheme returns the base color schema depending on the dark/light setting
func NewBaseTheme(preset Preset) *Theme {
	switch preset {
	case PresetDark:
		return &Theme{
			Default: color.MustParse("green"),
			Base: ThemeBase{
				Key:       color.MustParseSlice("hicyan / cyan"),
				Info:      color.MustParse("white"),
				Primary:   color.MustParse("magenta"),
				Secondary: color.MustParse("cyan"),
				Success:   color.MustParse("green"),
				Warning:   color.MustParse("yellow"),
				Danger:    color.MustParse("red"),
				Muted:     color.MustParse("gray:italic"),
			},
			Table: ThemeTable{
				Header: color.MustParse("bold"),
			},
			Data: ThemeData{
				String: color.MustParse("hiyellow"),
			},
		}

	case PresetLight:
		return &Theme{
			Default: color.MustParse("green"),
			Base: ThemeBase{
				Key:       color.MustParseSlice("hiblue / blue"),
				Info:      color.MustParse("black"),
				Primary:   color.MustParse("magenta"),
				Secondary: color.MustParse("blue"),
				Success:   color.MustParse("green"),
				Warning:   color.MustParse("yellow"),
				Danger:    color.MustParse("red"),
				Muted:     color.MustParse("gray:italic"),
			},
			Table: ThemeTable{
				Header: color.MustParse("bold"),
			},
			Data: ThemeData{
				String: color.MustParse("yellow"),
			},
		}

	// Special Preset for Protanopias
	case PresetProtDark:
		return &Theme{
			Default: color.MustParse("white"),
			Base: ThemeBase{
				Key:       color.MustParseSlice("#feb927 / #fe6e1a"),
				Info:      color.MustParse("white"),
				Primary:   color.MustParse("#4860e6"),             // magenta
				Secondary: color.MustParse("#2aabee"),             // cyan
				Success:   color.MustParse("#6afd6a:bold"),        // bold green
				Warning:   color.MustParse("#feb927:italic"),      // yellow
				Danger:    color.MustParse("fg=white:bg=#c2270a"), // red background
				Muted:     color.MustParse("#2ee5ae:italic"),      // white-ish
			},
			Data: ThemeData{
				String: color.MustParse("#2aabee"),
			},
			Table: ThemeTable{
				Header:  color.MustParse("white:bold"),
				Columns: color.MustParseSlice("#2aabee / #6afd6a:bold / #4860e6 / white / #feb927"),
			},
		}

	case PresetProtLight:
		return &Theme{
			Default: color.MustParse("black"),
			Base: ThemeBase{
				Key:       color.MustParseSlice("#feb927 / #fe6e1a"),
				Info:      color.MustParse("black"),
				Primary:   color.MustParse("#4860e6"),
				Secondary: color.MustParse("#2aabee"),
				Success:   color.MustParse("#6afd6a:bold"),
				Warning:   color.MustParse("#feb927:italic"),
				Danger:    color.MustParse("fg=black:bg=#c2270a"),
				Muted:     color.MustParse("#2ee5ae:italic"),
			},
			Data: ThemeData{
				String: color.MustParse("#2aabee"),
			},
			Table: ThemeTable{
				Header:  color.MustParse("black:bold"),
				Columns: color.MustParseSlice("#2aabee / #6afd6a:bold / #4860e6 / black / #feb927"),
			},
		}

	// Special Preset for Deuteranopia
	case PresetDeutDark:
		return &Theme{
			Default: color.MustParse("white"),
			Base: ThemeBase{
				Key:       color.MustParseSlice("#feb927 / #fe6e1a"),
				Info:      color.MustParse("white"),
				Primary:   color.MustParse("#4860e6"),
				Secondary: color.MustParse("#2aabee"),
				Success:   color.MustParse("#6afd6a:bold"),
				Warning:   color.MustParse("#feb927:italic"),
				Danger:    color.MustParse("fg=white:bg=#c2270a"),
				Muted:     color.MustParse("#2ee5ae"),
			},
			Data: ThemeData{
				String: color.MustParse("#2aabee"),
			},
			Table: ThemeTable{
				Header:  color.MustParse("white:bold"),
				Columns: color.MustParseSlice("#2aabee / #6afd6a:bold / #4860e6 / white / #feb927"),
			},
		}

	case PresetDeutLight:
		return &Theme{
			Default: color.MustParse("black"),
			Base: ThemeBase{
				Key:       color.MustParseSlice("#feb927 / #fe6e1a"),
				Info:      color.MustParse("black"),
				Primary:   color.MustParse("#4860e6"),
				Secondary: color.MustParse("#2aabee"),
				Success:   color.MustParse("#6afd6a:bold"),
				Warning:   color.MustParse("#feb927:italic"),
				Danger:    color.MustParse("fg=black:bg=#c2270a"),
				Muted:     color.MustParse("#2ee5ae:italic"),
			},
			Data: ThemeData{
				String: color.MustParse("#2aabee"),
			},
			Table: ThemeTable{
				Header:  color.MustParse("black:bold"),
				Columns: color.MustParseSlice("#2aabee / #6afd6a:bold / #4860e6 / black / #feb927"),
			},
		}

	// Special Preset for Tritanopia
	case PresetTritDark:
		return &Theme{
			Default: color.MustParse("white"),
			Base: ThemeBase{
				Key:       color.MustParseSlice("#feb927 / #fe6e1a"),
				Info:      color.MustParse("white"),
				Primary:   color.MustParse("#4860e6"),
				Secondary: color.MustParse("#2aabee"),
				Success:   color.MustParse("#6afd6a:bold"),
				Warning:   color.MustParse("#feb927:italic"),
				Danger:    color.MustParse("fg=white:bg=#c2270a"),
				Muted:     color.MustParse("#2ee5ae"),
			},
			Data: ThemeData{
				String: color.MustParse("#2aabee"),
			},
			Table: ThemeTable{
				Header:  color.MustParse("white:bold"),
				Columns: color.MustParseSlice("#2aabee / #6afd6a:bold / #4860e6 / white / #feb927"),
			},
		}

	case PresetTritLight:
		return &Theme{
			Default: color.MustParse("black"),
			Base: ThemeBase{
				Key:       color.MustParseSlice("#feb927 / #fe6e1a"),
				Info:      color.MustParse("black"),
				Primary:   color.MustParse("#4860e6"),
				Secondary: color.MustParse("#2aabee"),
				Success:   color.MustParse("#6afd6a:bold"),
				Warning:   color.MustParse("#feb927:italic"),
				Danger:    color.MustParse("fg=black:bg=#c2270a"),
				Muted:     color.MustParse("#2ee5ae:italic"),
			},
			Data: ThemeData{
				String: color.MustParse("#2aabee"),
			},
			Table: ThemeTable{
				Header:  color.MustParse("black:bold"),
				Columns: color.MustParseSlice("#2aabee / #6afd6a:bold / #4860e6 / black / #feb927"),
			},
		}

	// Pre-v0.3.0
	case PresetPre030Dark:
		return &Theme{
			Default: color.MustParse("green"),
			Base: ThemeBase{
				Key:       color.MustParseSlice("yellow / white"),
				Info:      color.MustParse("white"),
				Primary:   color.MustParse("magenta"),
				Secondary: color.MustParse("cyan"),
				Success:   color.MustParse("green"),
				Warning:   color.MustParse("yellow"),
				Danger:    color.MustParse("red"),
				Muted:     color.MustParse("yellow"),
			},
			Options: ThemeOptions{
				Flag: color.MustParse("yellow"),
			},
		}

	case PresetPre030Light:
		return &Theme{
			Default: color.MustParse("green"),
			Base: ThemeBase{
				Key:       color.MustParseSlice("yellow / black"),
				Info:      color.MustParse("black"),
				Primary:   color.MustParse("magenta"),
				Secondary: color.MustParse("blue"),
				Success:   color.MustParse("green"),
				Warning:   color.MustParse("yellow"),
				Danger:    color.MustParse("red"),
				Muted:     color.MustParse("yellow"),
			},
			Options: ThemeOptions{
				Flag: color.MustParse("yellow"),
			},
		}

	// Pre-v0.0.21
	case PresetPre0021Dark:
		return &Theme{
			Default: color.MustParse("green"),
			Base: ThemeBase{
				Key:       color.MustParseSlice("yellow / white"),
				Info:      color.MustParse("white"),
				Primary:   color.MustParse("magenta"),
				Secondary: color.MustParse("cyan"),
				Success:   color.MustParse("green"),
				Warning:   color.MustParse("yellow"),
				Danger:    color.MustParse("red"),
				Muted:     color.MustParse("yellow"),
			},
			Data: ThemeData{
				String: color.MustParse("cyan"),
			},
			Table: ThemeTable{
				Columns: color.MustParseSlice("cyan / green / magenta / white / yellow"),
			},
			Status: ThemeStatus{
				Success: color.MustParse("none"),
			},
			Options: ThemeOptions{
				Flag: color.MustParse("yellow"),
			},
		}

	case PresetPre0021Light:
		return &Theme{
			Default: color.MustParse("green"),
			Base: ThemeBase{
				Key:       color.MustParseSlice("yellow / black"),
				Info:      color.MustParse("white"),
				Primary:   color.MustParse("magenta"),
				Secondary: color.MustParse("cyan"),
				Success:   color.MustParse("green"),
				Warning:   color.MustParse("yellow"),
				Danger:    color.MustParse("red"),
				Muted:     color.MustParse("yellow"),
			},
			Data: ThemeData{
				String: color.MustParse("blue"),
			},
			Table: ThemeTable{
				Columns: color.MustParseSlice("cyan / green / magenta / black / yellow / blue"),
			},
			Status: ThemeStatus{
				Success: color.MustParse("none"),
			},
			Options: ThemeOptions{
				Flag: color.MustParse("yellow"),
			},
		}

	default:
		// Empty theme
		return &Theme{}
	}
}

// Theme is the root theme config.
type Theme struct {
	// Base colors must be first so they're applied first
	Base ThemeBase // base colors for themes

	Default color.Color // default when no specific mapping is found for the command

	Shell  ThemeShell  // colors for representing shells (e.g bash, zsh, etc)
	Data   ThemeData   // colors for representing data
	Status ThemeStatus // generic status coloring (e.g "Ready", "Terminating")
	Table  ThemeTable  // used in table output, e.g "kubectl get" and parts of "kubectl describe"
	Stderr ThemeStderr // used in kubectl's stderr output

	Apply    ThemeApply    // used in "kubectl apply"
	Create   ThemeCreate   // used in "kubectl create"
	Delete   ThemeDelete   // used in "kubectl delete"
	Describe ThemeDescribe // used in "kubectl describe"
	Drain    ThemeDrain    // used in "kubectl drain"
	Explain  ThemeExplain  // used in "kubectl explain"
	Expose   ThemeExpose   // used in "kubectl expose"
	Options  ThemeOptions  // used in "kubectl options"
	Patch    ThemePatch    // used in "kubectl patch"
	Rollout  ThemeRollout  // used in "kubectl rollout"
	Scale    ThemeScale    // used in "kubectl scale"
	Uncordon ThemeUncordon // used in "kubectl uncordon"
	Version  ThemeVersion  // used in "kubectl version"
	Help     ThemeHelp     // used in "kubectl --help"
	Logs     ThemeLogs     // used in "kubectl logs"
	Diff     ThemeDiff     // used in "kubectl diff"
}

func (t *Theme) ComputeCache() {
	themeVal := reflect.ValueOf(t).Elem()
	walkFields(themeVal, "theme", visitorComputeCache)
}

// ThemeBase contains base colors that other theme fields can default to,
// just to make overriding themes easier.
//
// These fields should never be referenced in the printers.
// Instead, they should use the more specific fields, such as [ThemeApply.Created]
type ThemeBase struct {
	Danger    color.Color // general color for when things are bad
	Info      color.Color // general color for when things are informational
	Muted     color.Color // general color for when things are less relevant
	Primary   color.Color // general color for when things are focus
	Secondary color.Color // general color for when things are secondary focus
	Success   color.Color // general color for when things are good
	Warning   color.Color // general color for when things are wrong

	Key color.Slice `defaultFromMany:"theme.base.secondary"` // general color for keys
}

// ThemeShell holds colors for when representing shell commands (bash, zsh, etc)
type ThemeShell struct {
	Comment color.Color `defaultFrom:"theme.base.muted"`     // used on comments, e.g `# this is a comment`
	Command color.Color `defaultFrom:"theme.base.success"`   // used on commands, e.g `kubectl` or `echo`
	Arg     color.Color `defaultFrom:"theme.base.info"`      // used on arguments, e.g `get pods` in `kubectl get pods`
	Flag    color.Color `defaultFrom:"theme.base.secondary"` // used on flags, e.g `--watch` in `kubectl get pods --watch`
}

// ThemeData holds colors for when representing parsed data.
// Such as in YAML, JSON, and even some "kubectl describe" values
type ThemeData struct {
	Key    color.Slice `defaultFrom:"theme.base.key"`     // used for the key
	String color.Color `defaultFrom:"theme.base.info"`    // used when value is a string
	True   color.Color `defaultFrom:"theme.base.success"` // used when value is true
	False  color.Color `defaultFrom:"theme.base.danger"`  // used when value is false
	Number color.Color `defaultFrom:"theme.base.primary"` // used when the value is a number
	Null   color.Color `defaultFrom:"theme.base.muted"`   // used when the value is null, nil, or none

	Quantity      color.Color `defaultFrom:"theme.data.number"`  // used when the value is a quantity, e.g "100m" or "5Gi"
	Duration      color.Color ``                                 // used when the value is a duration, e.g "12m" or "1d12h"
	DurationFresh color.Color `defaultFrom:"theme.base.success"` // color used when the time value is under a certain delay

	Ratio ThemeDataRatio
}

type ThemeDataRatio struct {
	Zero    color.Color `defaultFrom:"theme.base.muted"`   // used for "0/0"
	Equal   color.Color ``                                 // used for "n/n", e.g "1/1"
	Unequal color.Color `defaultFrom:"theme.base.warning"` // used for "n/m", e.g "0/1"
}

// ThemeStatus holds colors for status texts, used in for example
// the "kubectl get" status column
type ThemeStatus struct {
	Success color.Color `defaultFrom:"theme.base.success"` // used in status keywords, e.g "Running", "Ready"
	Warning color.Color `defaultFrom:"theme.base.warning"` // used in status keywords, e.g "Terminating"
	Error   color.Color `defaultFrom:"theme.base.danger"`  // used in status keywords, e.g "Failed", "Unhealthy"
}

// ThemeTable holds colors for table output
type ThemeTable struct {
	Header  color.Color `defaultFrom:"theme.base.info"`                          // used on table headers
	Columns color.Slice `defaultFromMany:"theme.base.info,theme.base.secondary"` // used on table columns when no other coloring applies such as status or duration coloring. The multiple colors are cycled based on column ID, from left to right.
}

// ThemeStderr holds generic colors for kubectl's stderr output.
type ThemeStderr struct {
	Error color.Color `defaultFrom:"theme.base.danger"` // e.g when text contains "error"

	NoneFound          color.Color `defaultFrom:"theme.data.null"`   // used on table output like "No resources found"
	NoneFoundNamespace color.Color `defaultFrom:"theme.data.string"` // used on the namespace name of "No resources found in my-ns namespace"

	// default when no specific mapping is found for the output line
	//
	// Deprecated: This field is no longer used (since v0.4.0),
	// as the stderr logs now uses the "kubectl logs" behavior as a fallback/default coloring.
	Default color.Color `jsonschema_extras:"deprecated=true"` // *deprecated: this field is no longer used (since v0.4.0)*
}

// ThemeApply holds colors for the "kubectl apply" output.
type ThemeDescribe struct {
	Key color.Slice `defaultFrom:"theme.base.key"` // used on keys. The multiple colors are cycled based on indentation.
}

// ThemeApply holds colors for the "kubectl apply" output.
type ThemeApply struct {
	Created    color.Color `defaultFrom:"theme.base.success"` // used on "deployment.apps/foo created"
	Configured color.Color `defaultFrom:"theme.base.warning"` // used on "deployment.apps/bar configured"
	Unchanged  color.Color `defaultFrom:"theme.base.primary"` // used on "deployment.apps/quux unchanged"
	Serverside color.Color `defaultFrom:"theme.base.warning"` // used on "deployment.apps/quux serverside-applied"

	DryRun   color.Color `defaultFrom:"theme.base.secondary"` // used on "(dry run)" and "(server dry run)"
	Fallback color.Color `defaultFrom:"theme.base.success"`   // used when outputs unknown format
}

// ThemeDelete holds colors for the "kubectl delete" output.
type ThemeDelete struct {
	Deleted color.Color `defaultFrom:"theme.base.danger"` // used on "deployment.apps "nginx" deleted"

	DryRun   color.Color `defaultFrom:"theme.apply.dryrun"` // used on "(dry run)" and "(server dry run)"
	Fallback color.Color `defaultFrom:"theme.base.danger"`  // used when outputs unknown format
}

// ThemeCreate holds colors for the "kubectl create" output.
type ThemeCreate struct {
	Created color.Color `defaultFrom:"theme.base.success"` // used on "deployment.apps/foo created"

	DryRun   color.Color `defaultFrom:"theme.apply.dryrun"` // used on "(dry run)" and "(server dry run)"
	Fallback color.Color `defaultFrom:"theme.base.success"` // used when outputs unknown format
}

// ThemeExpose holds colors for the "kubectl expose" output.
type ThemeExpose struct {
	Exposed color.Color `defaultFrom:"theme.base.primary"` // used on "deployment.apps/foo created"

	DryRun   color.Color `defaultFrom:"theme.apply.dryrun"` // used on "(dry run)" and "(server dry run)"
	Fallback color.Color `defaultFrom:"theme.base.primary"` // used when outputs unknown format
}

// ThemeScale holds colors for the "kubectl scale" output.
type ThemeScale struct {
	Scaled color.Color `defaultFrom:"theme.base.warning"` // used on "deployment.apps/foo scaled"

	DryRun   color.Color `defaultFrom:"theme.apply.dryrun"` // used on "(dry run)" and "(server dry run)"
	Fallback color.Color `defaultFrom:"theme.base.warning"` // used when outputs unknown format
}

// ThemeRollout holds colors for the "kubectl rollout" output.
type ThemeRollout struct {
	RolledBack color.Color `defaultFrom:"theme.base.warning"`   // used on "deployment.apps/foo rolled back"
	Paused     color.Color `defaultFrom:"theme.base.primary"`   // used on "deployment.apps/foo paused"
	Resumed    color.Color `defaultFrom:"theme.base.secondary"` // used on "deployment.apps/foo resumed"
	Restarted  color.Color `defaultFrom:"theme.base.warning"`   // used on "deployment.apps/foo restarted"

	DryRun   color.Color `defaultFrom:"theme.apply.dryrun"` // used on "(dry run)" and "(server dry run)"
	Fallback color.Color `defaultFrom:"theme.base.warning"` // used when outputs unknown format
}

// ThemePatch holds colors for the "kubectl patch" output.
type ThemePatch struct {
	Patched color.Color `defaultFrom:"theme.base.warning"` // used on "deployment.apps/foo patched"

	DryRun   color.Color `defaultFrom:"theme.apply.dryrun"` // used on "(dry run)" and "(server dry run)"
	Fallback color.Color `defaultFrom:"theme.base.warning"` // used when outputs unknown format
}

// ThemeUncordon holds colors for the "kubectl uncordon" output.
type ThemeUncordon struct {
	Uncordoned color.Color `defaultFrom:"theme.base.secondary"` // used on "node/my-worker-node-01 uncordoned"

	DryRun   color.Color `defaultFrom:"theme.apply.dryrun"` // used on "(dry run)" and "(server dry run)"
	Fallback color.Color `defaultFrom:"theme.base.warning"` // used when outputs unknown format
}

// ThemeDrain holds colors for the "kubectl drain" output.
type ThemeDrain struct {
	Cordoned    color.Color `defaultFrom:"theme.base.primary"` // used on "node/my-worker-node-01 cordoned"
	EvictingPod color.Color `defaultFrom:"theme.base.muted"`   // used on "evicting pod my-namespace/my-pod"
	Evicted     color.Color `defaultFrom:"theme.base.warning"` // used on "pod/my-pod evicted"
	Drained     color.Color `defaultFrom:"theme.base.success"` // used on "node/my-worker-node-01 drained"

	DryRun   color.Color `defaultFrom:"theme.apply.dryrun"` // used on "(dry run)" and "(server dry run)"
	Fallback color.Color `defaultFrom:"theme.base.warning"` // used when outputs unknown format
}

// ThemeExplain holds colors for the "kubectl explain" output.
type ThemeExplain struct {
	Key      color.Slice `defaultFrom:"theme.base.key"`    // used on keys. The multiple colors are cycled based on indentation.
	Required color.Color `defaultFrom:"theme.base.danger"` // used on the trailing "-required-" string
}

// ThemeDiff holds colors for the "kubectl diff" output.
type ThemeDiff struct {
	Added     color.Color `defaultFrom:"theme.base.success"` // used on added lines
	Removed   color.Color `defaultFrom:"theme.base.danger"`  // used on removed lines
	Unchanged color.Color `defaultFrom:"theme.base.muted"`   // used on unchanged lines
}

// ThemeOptions holds colors for the "kubectl options" output.
type ThemeOptions struct {
	Flag color.Color `defaultFrom:"theme.base.secondary"` // e.g "--kubeconfig"
}

// ThemeVersion holds colors for the "kubectl version" output.
type ThemeVersion struct {
	Key color.Slice `defaultFrom:"theme.base.key"` // used on the key
}

// ThemeHelp holds colors for the "kubectl --help" output.
type ThemeHelp struct {
	Header   color.Color `defaultFrom:"theme.table.header"`   // e.g "Examples:" or "Options:"
	Flag     color.Color `defaultFrom:"theme.base.secondary"` // e.g "--kubeconfig"
	FlagDesc color.Color `defaultFrom:"theme.base.info"`      // Flag descripion under "Options:" heading
	Url      color.Color `defaultFrom:"theme.base.secondary"` // e.g `[https://example.com]`
	Text     color.Color `defaultFrom:"theme.base.info"`      // Fallback text color
}

// ThemeLogs holds colors for the "kubectl logs" output.
type ThemeLogs struct {
	Key          color.Slice `defaultFrom:"theme.data.key"`
	QuotedString color.Color `defaultFrom:"theme.data.string"` // Used on quoted strings that are not part of a `key="value"`
	Date         color.Color `defaultFrom:"theme.base.muted"`
	SourceRef    color.Color `defaultFrom:"theme.base.muted"`
	GUID         color.Color `defaultFrom:"theme.base.muted"`

	Severity ThemeLogsSeverity
}

// ThemeLogsSeverity holds colors for "log level severity" found in "kubectl logs" output
type ThemeLogsSeverity struct {
	Trace color.Color `defaultFrom:"theme.base.muted"`
	Debug color.Color `defaultFrom:"theme.base.muted"`
	Info  color.Color `defaultFrom:"theme.base.success"`
	Warn  color.Color `defaultFrom:"theme.base.warning"`
	Error color.Color `defaultFrom:"theme.base.danger"`
	Fatal color.Color `defaultFrom:"theme.base.danger"`
	Panic color.Color `defaultFrom:"theme.base.danger"`
}

func applyViperDefaults(theme *Theme, v *viper.Viper) {
	themeVal := reflect.ValueOf(theme).Elem()
	walkFields(themeVal, "theme", themeViperVisitor{viper: v}.visitorApplyDefaults)
}

type themeViperVisitor struct {
	viper *viper.Viper
}

func walkFields(val reflect.Value, viperKey string, visitor func(viperKey string, value reflect.Value, tags reflect.StructTag)) {
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
			walkFields(fieldVal, newViperKey, visitor)
			continue
		}
		visitor(newViperKey, fieldVal, fieldTyp.Tag)
	}
}

func visitorComputeCache(viperKey string, value reflect.Value, _ reflect.StructTag) {
	switch value := value.Addr().Interface().(type) {
	case *color.Color:
		value.ComputeCache()
	case *color.Slice:
		value.ComputeCache()
	default:
		panic(fmt.Errorf("%s: unsupported field type: %T", viperKey, value))
	}
}

func (t themeViperVisitor) visitorApplyDefaults(viperKey string, value reflect.Value, tags reflect.StructTag) {
	switch value := value.Interface().(type) {
	case color.Color:
		if _, ok := tags.Lookup("defaultFromMany"); ok {
			panic(fmt.Errorf("%s: cannot use defaultFromMany tag on a Color field", viperKey))
		}
		if defaultFrom, ok := tags.Lookup("defaultFrom"); ok {
			t.setColorOrKey(viperKey, value, defaultFrom)
		} else {
			t.setColor(viperKey, value)
		}
	case color.Slice:
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

func (t themeViperVisitor) setColorOrKey(key string, value color.Color, otherKey string) {
	if t.setColor(key, value) {
		return
	}
	t.viper.SetDefault(key, t.viper.Get(otherKey))
}

func (t themeViperVisitor) setColor(key string, value color.Color) bool {
	if value.IsZero() {
		t.viper.SetDefault(key, color.Color{})
		return false
	}
	t.viper.SetDefault(key, value)
	return true
}

func (t themeViperVisitor) setColorSliceOrKey(key string, value color.Slice, otherKey string) {
	if t.setColorSlice(key, value) {
		return
	}
	t.viper.SetDefault(key, t.viper.Get(otherKey))
}

func (t themeViperVisitor) setColorSliceOrManyKeys(key string, value color.Slice, otherKeys []string) {
	if t.setColorSlice(key, value) {
		return
	}
	values := make(color.Slice, 0, len(otherKeys))
	for _, k := range otherKeys {
		val := t.viper.Get(k)
		col, ok := val.(color.Color)
		if !ok {
			col = color.MustParse(fmt.Sprint(val))
		}
		if val != nil {
			values = append(values, col)
		}
	}
	if len(values) > 0 {
		t.viper.SetDefault(key, values)
	}
}

func (t themeViperVisitor) setColorSlice(key string, value color.Slice) bool {
	if len(value) == 0 {
		t.viper.SetDefault(key, color.Slice{})
		return false
	}
	t.viper.SetDefault(key, value)
	return true
}
