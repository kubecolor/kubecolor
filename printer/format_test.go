package printer

import (
	"testing"

	"github.com/kubecolor/kubecolor/color"
)

func Test_toSpaces(t *testing.T) {
	if toSpaces(3) != "   " {
		t.Fatalf("fail")
	}
}

func Test_getColorByKeyIndent(t *testing.T) {
	tests := []struct {
		name             string
		themePreset      color.Preset
		indent           int
		basicIndentWidth int
		expected         color.Color
	}{
		{"dark depth: 1", color.PresetDark, 2, 2, color.White},
		{"light depth: 1", color.PresetLight, 2, 2, color.Black},
		{"dark depth: 2", color.PresetDark, 4, 2, color.Yellow},
		{"light depth: 2", color.PresetLight, 4, 2, color.Yellow},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getColorByKeyIndent(tt.indent, tt.basicIndentWidth, color.NewTheme(tt.themePreset))
			if got != tt.expected {
				t.Errorf("fail: got: %v, expected: %v", got, tt.expected)
			}
		})
	}
}

func Test_getColorByValueType(t *testing.T) {
	tests := []struct {
		name        string
		themePreset color.Preset
		val         string
		expected    color.Color
	}{
		{"dark null", color.PresetDark, "null", color.Yellow},
		{"light null", color.PresetLight, "<none>", color.Yellow},

		{"dark true", color.PresetDark, "true", color.Green},
		{"light true", color.PresetLight, "true", color.Green},

		{"dark false", color.PresetDark, "false", color.Red},
		{"light false", color.PresetLight, "false", color.Red},

		{"dark number", color.PresetDark, "123", color.Magenta},
		{"light number", color.PresetLight, "456", color.Magenta},

		{"dark string", color.PresetDark, "aaa", color.White},
		{"light string", color.PresetLight, "12345a", color.Black},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getColorByValueType(tt.val, color.NewTheme(tt.themePreset))
			if got != tt.expected {
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
