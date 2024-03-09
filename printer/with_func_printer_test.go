package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_WithFuncPrinter_Print(t *testing.T) {
	var (
		colorWhite = config.MustParseColor("white")
		colorRed   = config.MustParseColor("red")
	)
	tests := []struct {
		name     string
		fn       func(line string) config.Color
		input    string
		expected string
	}{
		{
			name: "colored in white",
			fn: func(_ string) config.Color {
				return colorWhite
			},
			input: testutil.NewHereDoc(`
				test
				test2
				test3`),
			expected: testutil.NewHereDocf(`
				%s
				%s
				%s
				`, colorWhite.Render("test"), colorWhite.Render("test2"), colorWhite.Render("test3")),
		},
		{
			name: "color changes by line",
			fn: func(line string) config.Color {
				if line == "test2" {
					return colorRed
				}
				return colorWhite
			},
			input: testutil.NewHereDoc(`
				test
				test2
				test3`),
			expected: testutil.NewHereDocf(`
				%s
				%s
				%s
				`, colorWhite.Render("test"), colorRed.Render("test2"), colorWhite.Render("test3")),
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
