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
		dark             bool
		indent           int
		basicIndentWidth int
		expected         color.Color
	}{
		{"dark depth: 1", true, 2, 2, color.White},
		{"light depth: 1", false, 2, 2, color.Black},
		{"dark depth: 2", true, 4, 2, color.Yellow},
		{"light depth: 2", false, 4, 2, color.Yellow},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getColorByKeyIndent(tt.indent, tt.basicIndentWidth, NewColorSchema(tt.dark))
			if got != tt.expected {
				t.Errorf("fail: got: %v, expected: %v", got, tt.expected)
			}
		})
	}
}

func Test_getColorByValueType(t *testing.T) {
	tests := []struct {
		name     string
		dark     bool
		val      string
		expected color.Color
	}{
		{"dark null", true, "null", NullColorForDark},
		{"light null", false, "<none>", NullColorForLight},

		{"dark true", true, "true", TrueColorForDark},
		{"light true", false, "true", TrueColorForLight},

		{"dark false", true, "false", FalseColorForDark},
		{"light false", false, "false", FalseColorForLight},

		{"dark number", true, "123", NumberColorForDark},
		{"light number", false, "456", NumberColorForLight},

		{"dark string", true, "aaa", StringColorForDark},
		{"light string", false, "12345a", StringColorForLight},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getColorByValueType(tt.val, NewColorSchema(tt.dark))
			if got != tt.expected {
				t.Errorf("fail: got: %v, expected: %v", got, tt.expected)
			}
		})
	}
}

// func Test_getColorsByBackground(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		dark     bool
// 		expected []color.Color
// 	}{
// 		{"dark", true, colorsForDarkBackground},
// 		{"light", false, colorsForLightBackground},
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			got := getColorsByBackground(tt.dark)
// 			if diff := cmp.Diff(got, tt.expected); diff != "" {
// 				t.Errorf("fail: %v", diff)
// 			}
// 		})
// 	}
// }

// func Test_getHeaderColorByBackground(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		dark     bool
// 		expected color.Color
// 	}{
// 		{"dark", true, HeaderColorForDark},
// 		{"light", false, HeaderColorForLight},
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			got := getHeaderColorByBackground(tt.dark)
// 			if got != tt.expected {
// 				t.Errorf("fail: got: %v, expected: %v", got, tt.expected)
// 			}
// 		})
// 	}
// }

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
