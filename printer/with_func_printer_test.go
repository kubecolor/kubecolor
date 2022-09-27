package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/color"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_WithFuncPrinter_Print(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(line string) color.Color
		input    string
		expected string
	}{
		{
			name: "colored in white",
			fn: func(_ string) color.Color {
				return color.White
			},
			input: testutil.NewHereDoc(`
				test
				test2
				test3`),
			expected: testutil.NewHereDocf(`
				%s
				%s
				%s
				`, color.Apply("test", color.White), color.Apply("test2", color.White), color.Apply("test3", color.White)),
		},
		{
			name: "color changes by line",
			fn: func(line string) color.Color {
				if line == "test2" {
					return color.Red
				}
				return color.White
			},
			input: testutil.NewHereDoc(`
				test
				test2
				test3`),
			expected: testutil.NewHereDocf(`
				%s
				%s
				%s
				`, color.Apply("test", color.White), color.Apply("test2", color.Red), color.Apply("test3", color.White)),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := WithFuncPrinter{Fn: tt.fn}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
