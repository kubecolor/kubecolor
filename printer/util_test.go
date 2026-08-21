package printer

import (
	"testing"
	"time"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/color"
	"github.com/kubecolor/kubecolor/config/testconfig"
)

func Test_ColorDuration(t *testing.T) {
	// Distinct colors so overflow (data.duration) is never confused with a
	// gradient color (data.durationFresh).
	theme := &config.Theme{
		Data: config.ThemeData{
			DurationFresh: color.MustParseSlice("green/yellow/red"),
			Duration:      color.MustParse("cyan"), // overflow / fallback color
		},
	}

	tests := []struct {
		name       string
		thresholds config.DurationSlice
		age        time.Duration
		expected   string
	}{
		{"empty thresholds -> overflow", config.DurationSlice{}, 30 * time.Minute, "cyan"},

		{"single threshold, below -> first color", config.MustParseDurationSlice("1h"), 30 * time.Minute, "green"},
		{"single threshold, at/above -> overflow", config.MustParseDurationSlice("1h"), 2 * time.Hour, "cyan"},

		{"three thresholds, bucket 0", config.MustParseDurationSlice("1h/6h/1d"), 30 * time.Minute, "green"},
		{"three thresholds, bucket 1", config.MustParseDurationSlice("1h/6h/1d"), 3 * time.Hour, "yellow"},
		{"three thresholds, bucket 2", config.MustParseDurationSlice("1h/6h/1d"), 12 * time.Hour, "red"},
		{"three thresholds, overflow", config.MustParseDurationSlice("1h/6h/1d"), 48 * time.Hour, "cyan"},

		{"more thresholds than colors -> surplus bucket uses overflow", config.MustParseDurationSlice("1h/6h/1d/7d"), 48 * time.Hour, "cyan"},

		{"more colors than thresholds -> extras ignored, overflow above", config.MustParseDurationSlice("1h"), 2 * time.Hour, "cyan"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ColorDuration(tt.age, tt.thresholds, theme)
			if got.String() != tt.expected {
				t.Errorf("fail: got: %q, expected: %q", got, tt.expected)
			}
		})
	}
}

// When theme.data.duration is unset, overflow ages must behave differently
// depending on whether thresholds are configured:
//   - no thresholds  -> noop color, so table column cycling still applies
//   - thresholds set -> the first table column color, so it reads as a normal
//     column instead of the arbitrary cycled color for the age column's index
func Test_ColorDuration_unsetOverflow(t *testing.T) {
	theme := &config.Theme{
		Table: config.ThemeTable{
			Columns: color.MustParseSlice("white/cyan/magenta"),
		},
		Data: config.ThemeData{
			DurationFresh: color.MustParseSlice("green/yellow"),
			// Duration intentionally left unset (noop)
		},
	}

	noThresholds := ColorDuration(48*time.Hour, config.DurationSlice{}, theme)
	if !noThresholds.IsNoop() {
		t.Errorf("no thresholds: expected noop color (allow cycling), got %q", noThresholds)
	}

	withThresholds := ColorDuration(48*time.Hour, config.MustParseDurationSlice("1h/24h"), theme)
	if withThresholds.String() != "white" {
		t.Errorf("with thresholds: expected first table column color, got %q", withThresholds)
	}
}

func Test_getColorByKeyIndent(t *testing.T) {
	tests := []struct {
		name             string
		theme            *config.Theme
		indent           int
		basicIndentWidth int
		expected         string
	}{
		{"dark depth: 1", testconfig.DarkTheme, 2, 2, "cyan"},
		{"light depth: 1", testconfig.LightTheme, 2, 2, "blue"},
		{"dark depth: 2", testconfig.DarkTheme, 4, 2, "hicyan"},
		{"light depth: 2", testconfig.LightTheme, 4, 2, "hiblue"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ColorDataKey(tt.indent, tt.basicIndentWidth, tt.theme.Base.Key)
			if got.String() != tt.expected {
				t.Errorf("fail: got: %q, expected: %q", got, tt.expected)
			}
		})
	}
}

func Test_getColorByValueType(t *testing.T) {
	tests := []struct {
		name     string
		theme    *config.Theme
		val      string
		expected string
	}{
		{"dark null", testconfig.DarkTheme, "null", "gray:italic"},
		{"light null", testconfig.LightTheme, "<none>", "gray:italic"},

		{"dark true", testconfig.DarkTheme, "true", "green"},
		{"light true", testconfig.LightTheme, "true", "green"},

		{"dark false", testconfig.DarkTheme, "false", "red"},
		{"light false", testconfig.LightTheme, "false", "red"},

		{"dark number", testconfig.DarkTheme, "123", "magenta"},
		{"light number", testconfig.LightTheme, "456", "magenta"},

		{"dark string", testconfig.DarkTheme, "aaa", "hiyellow"},
		{"light string", testconfig.LightTheme, "12345a", "yellow"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ColorDataValue(tt.val, tt.theme)
			if got.String() != tt.expected {
				t.Errorf("fail: got: %v, expected: %v", got, tt.expected)
			}
		})
	}
}

func Test_findIndent(t *testing.T) {
	tests := []struct {
		line     string
		expected int
	}{
		{"no indent", 0},
		{"  2 indent", 2},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.line, func(t *testing.T) {
			t.Parallel()
			got := findIndent(tt.line)
			if got != tt.expected {
				t.Errorf("fail: got: %v, expected: %v", got, tt.expected)
			}
		})
	}
}
