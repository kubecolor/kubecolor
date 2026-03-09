package printer

import (
	"testing"
	"time"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/color"
	"github.com/kubecolor/kubecolor/config/testconfig"
)

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

func TestColorDuration(t *testing.T) {
	theme := testconfig.DarkTheme
	durationConfig := &config.Duration{
		Threshold1: 5 * time.Minute,
		Threshold2: 2 * time.Hour,
		Threshold3: 24 * time.Hour,
		Threshold4: 30 * 24 * time.Hour,
		Threshold5: 365 * 24 * time.Hour,
	}
	colors := &theme.Data.Duration

	tests := []struct {
		name   string
		column string
		wantOK bool
		want   color.Color
	}{
		{"zero duration", "0s", true, colors.Default},
		{"below threshold1", "30s", true, colors.Default},
		{"just under threshold1", "4m59s", true, colors.Default},
		{"at threshold1", "5m", true, colors.Threshold1},
		{"above threshold1", "10m", true, colors.Threshold1},
		{"at threshold2", "2h", true, colors.Threshold2},
		{"above threshold2", "3h", true, colors.Threshold2},
		{"at threshold3", "24h", true, colors.Threshold3},
		{"above threshold3", "5d", true, colors.Threshold3},
		{"at threshold4", "30d", true, colors.Threshold4},
		{"above threshold4", "45d", true, colors.Threshold4},
		{"at threshold5", "365d", true, colors.Threshold5},
		{"above threshold5", "2y", true, colors.Threshold5},
		{"invalid: not a duration", "Running", false, color.Color{}},
		{"invalid: empty string", "", false, color.Color{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ColorDuration(tt.column, durationConfig, theme)
			if ok != tt.wantOK {
				t.Fatalf("ColorDuration(%q) ok = %v, want %v", tt.column, ok, tt.wantOK)
			}
			if !ok {
				if got != tt.column {
					t.Errorf("ColorDuration(%q) on failure returned %q, want original string", tt.column, got)
				}
				return
			}
			if got != tt.want.Render(tt.column) {
				t.Errorf("ColorDuration(%q) = %q, want %q", tt.column, got, tt.want.Render(tt.column))
			}
		})
	}

	t.Run("zero thresholds: all durations use Default", func(t *testing.T) {
		zeroConfig := &config.Duration{}
		got, ok := ColorDuration("30s", zeroConfig, theme)
		if !ok {
			t.Fatal("ColorDuration(\"30s\") ok = false, want true")
		}
		if want := colors.Default.Render("30s"); got != want {
			t.Errorf("ColorDuration(\"30s\") = %q, want %q (Default)", got, want)
		}
	})

	t.Run("flat duration: all thresholds same color", func(t *testing.T) {
		flatColor := color.MustParse("yellow")
		flatTheme := testconfig.NewTheme(config.PresetDark)
		flatTheme.Data.Duration = config.ThemeDataDuration{
			Default:    flatColor,
			Threshold1: flatColor, Threshold2: flatColor, Threshold3: flatColor,
			Threshold4: flatColor, Threshold5: flatColor, Threshold6: flatColor,
		}

		for _, column := range []string{"10m", "3h", "5d", "2y"} {
			got, ok := ColorDuration(column, durationConfig, flatTheme)
			if !ok {
				t.Fatalf("ColorDuration(%q) ok = false, want true", column)
			}
			if want := flatColor.Render(column); got != want {
				t.Errorf("ColorDuration(%q) = %q, want %q", column, got, want)
			}
		}
	})
}
