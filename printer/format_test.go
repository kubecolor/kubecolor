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
		{"dark null", color.PresetDark, "null", NullColorForDark},
		{"light null", color.PresetLight, "<none>", NullColorForLight},

		{"dark true", color.PresetDark, "true", TrueColorForDark},
		{"light true", color.PresetLight, "true", TrueColorForLight},

		{"dark false", color.PresetDark, "false", FalseColorForDark},
		{"light false", color.PresetLight, "false", FalseColorForLight},

		{"dark number", color.PresetDark, "123", NumberColorForDark},
		{"light number", color.PresetLight, "456", NumberColorForLight},

		{"dark string", color.PresetDark, "aaa", StringColorForDark},
		{"light string", color.PresetLight, "12345a", StringColorForLight},
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
