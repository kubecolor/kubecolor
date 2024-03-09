package printer

import (
	"testing"

	"github.com/kubecolor/kubecolor/config"
)

func Test_toSpaces(t *testing.T) {
	if toSpaces(3) != "   " {
		t.Fatalf("fail")
	}
}

func Test_getColorByKeyIndent(t *testing.T) {
	tests := []struct {
		name             string
		themePreset      config.Preset
		indent           int
		basicIndentWidth int
		expected         string
	}{
		{"dark depth: 1", config.PresetDark, 2, 2, "white"},
		{"light depth: 1", config.PresetLight, 2, 2, "black"},
		{"dark depth: 2", config.PresetDark, 4, 2, "yellow"},
		{"light depth: 2", config.PresetLight, 4, 2, "yellow"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getColorByKeyIndent(tt.indent, tt.basicIndentWidth, config.NewTheme(tt.themePreset))
			if got.String() != tt.expected {
				t.Errorf("fail: got: %v, expected: %v", got, tt.expected)
			}
		})
	}
}

func Test_getColorByValueType(t *testing.T) {
	tests := []struct {
		name        string
		themePreset config.Preset
		val         string
		expected    string
	}{
		{"dark null", config.PresetDark, "null", "yellow"},
		{"light null", config.PresetLight, "<none>", "yellow"},

		{"dark true", config.PresetDark, "true", "green"},
		{"light true", config.PresetLight, "true", "green"},

		{"dark false", config.PresetDark, "false", "red"},
		{"light false", config.PresetLight, "false", "red"},

		{"dark number", config.PresetDark, "123", "magenta"},
		{"light number", config.PresetLight, "456", "magenta"},

		{"dark string", config.PresetDark, "aaa", "white"},
		{"light string", config.PresetLight, "12345a", "black"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getColorByValueType(tt.val, config.NewTheme(tt.themePreset))
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
