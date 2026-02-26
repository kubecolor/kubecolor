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
	colors := &theme.Data.DurationColors
	freshThreshold := 5 * time.Minute

	tests := []struct {
		name   string
		column string
		wantOK bool
		want   color.Color
	}{
		{
			name:   "fresh: under threshold",
			column: "30s",
			wantOK: true,
			want:   theme.Data.DurationFresh,
		},
		{
			name:   "fresh: exactly at threshold boundary",
			column: "4m59s",
			wantOK: true,
			want:   theme.Data.DurationFresh,
		},
		{
			name:   "boundary: exactly at threshold is NOT fresh",
			column: "5m",
			wantOK: true,
			want:   colors.Minutes,
		},
		{
			name:   "age-based: minutes",
			column: "10m",
			wantOK: true,
			want:   colors.Minutes,
		},
		{
			name:   "age-based: hours",
			column: "3h",
			wantOK: true,
			want:   colors.Hours,
		},
		{
			name:   "age-based: days",
			column: "5d",
			wantOK: true,
			want:   colors.Days,
		},
		{
			name:   "age-based: years",
			column: "2y",
			wantOK: true,
			want:   colors.Years,
		},
		{
			name:   "invalid: not a duration",
			column: "Running",
			wantOK: false,
		},
		{
			name:   "invalid: empty string",
			column: "",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ColorDuration(tt.column, freshThreshold, theme)
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

	t.Run("zero fresh threshold uses age-based seconds color", func(t *testing.T) {
		got, ok := ColorDuration("30s", 0, theme)
		if !ok {
			t.Fatal("ColorDuration(\"30s\") ok = false, want true")
		}
		if want := colors.Seconds.Render("30s"); got != want {
			t.Errorf("ColorDuration(\"30s\") = %q, want %q (Seconds)", got, want)
		}
	})

	t.Run("flat duration color overrides age-based coloring", func(t *testing.T) {
		flatColor := color.MustParse("yellow")
		flatTheme := testconfig.NewTheme(config.PresetDark)
		flatTheme.Data.DurationFlat = flatColor

		for _, column := range []string{"10m", "3h", "5d", "2y"} {
			got, ok := ColorDuration(column, freshThreshold, flatTheme)
			if !ok {
				t.Fatalf("ColorDuration(%q) ok = false, want true", column)
			}
			if want := flatColor.Render(column); got != want {
				t.Errorf("ColorDuration(%q) = %q, want %q", column, got, want)
			}
		}
	})

	t.Run("DurationFresh still applies when flat duration is set", func(t *testing.T) {
		flatTheme := testconfig.NewTheme(config.PresetDark)
		flatTheme.Data.DurationFlat = color.MustParse("yellow")

		got, ok := ColorDuration("30s", freshThreshold, flatTheme)
		if !ok {
			t.Fatal("ColorDuration(\"30s\") ok = false, want true")
		}
		if want := flatTheme.Data.DurationFresh.Render("30s"); got != want {
			t.Errorf("ColorDuration(\"30s\") = %q, want %q (DurationFresh)", got, want)
		}
	})
}

func Test_durationColorByAge(t *testing.T) {
	colors := &testconfig.DarkTheme.Data.DurationColors
	tests := []struct {
		name     string
		duration time.Duration
		want     color.Color
	}{
		{
			name:     "Empty duration",
			duration: 0,
			want:     colors.Seconds,
		},
		{
			name:     "30 seconds",
			duration: 30 * time.Second,
			want:     colors.Seconds,
		},
		{
			name:     "5 minutes",
			duration: 5 * time.Minute,
			want:     colors.Minutes,
		},
		{
			name:     "2 hours",
			duration: 2 * time.Hour,
			want:     colors.Hours,
		},
		{
			name:     "2 days",
			duration: 2 * 24 * time.Hour,
			want:     colors.Days,
		},
		{
			name:     "2 years",
			duration: 2 * 365 * 24 * time.Hour,
			want:     colors.Years,
		},
		{
			name:     "boundary: exactly 1 minute",
			duration: time.Minute,
			want:     colors.Minutes,
		},
		{
			name:     "boundary: exactly 1 hour",
			duration: time.Hour,
			want:     colors.Hours,
		},
		{
			name:     "boundary: exactly 24 hours",
			duration: 24 * time.Hour,
			want:     colors.Days,
		},
		{
			name:     "boundary: exactly 365 days",
			duration: 365 * 24 * time.Hour,
			want:     colors.Years,
		},
		{
			name:     "just under 1 minute",
			duration: 59 * time.Second,
			want:     colors.Seconds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := durationColorByAge(tt.duration, colors)

			if got.Render("x") != tt.want.Render("x") {
				t.Errorf("DurationColorByAge() color mismatch\nduration: %s\nexpected: %q\ngot:      %q",
					tt.duration, tt.want.Render("x"), got.Render("x"))
			}
		})
	}
}
