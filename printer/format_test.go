package printer

import (
	"testing"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/testconfig"
)

func Test_toSpaces(t *testing.T) {
	if toSpaces(3) != "   " {
		t.Fatalf("fail")
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
		{"dark depth: 1", testconfig.DarkTheme, 2, 2, "white"},
		{"light depth: 1", testconfig.LightTheme, 2, 2, "black"},
		{"dark depth: 2", testconfig.DarkTheme, 4, 2, "yellow"},
		{"light depth: 2", testconfig.LightTheme, 4, 2, "yellow"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getColorByKeyIndent(tt.indent, tt.basicIndentWidth, tt.theme)
			if got.String() != tt.expected {
				t.Errorf("fail: got: %v, expected: %v", got, tt.expected)
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
		{"dark null", testconfig.DarkTheme, "null", "yellow"},
		{"light null", testconfig.LightTheme, "<none>", "yellow"},

		{"dark true", testconfig.DarkTheme, "true", "green"},
		{"light true", testconfig.LightTheme, "true", "green"},

		{"dark false", testconfig.DarkTheme, "false", "red"},
		{"light false", testconfig.LightTheme, "false", "red"},

		{"dark number", testconfig.DarkTheme, "123", "magenta"},
		{"light number", testconfig.LightTheme, "456", "magenta"},

		{"dark string", testconfig.DarkTheme, "aaa", "white"},
		{"light string", testconfig.LightTheme, "12345a", "black"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getColorByValueType(tt.val, tt.theme)
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
