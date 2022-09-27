package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/color"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_SingleColoredPrinter_Print(t *testing.T) {
	tests := []struct {
		name     string
		color    color.Color
		input    string
		expected string
	}{
		{
			name:  "colored in white",
			color: color.White,
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
			name:  "colored in red",
			color: color.Red,
			input: testutil.NewHereDoc(`
				test
				test2
				test3`),
			expected: testutil.NewHereDocf(`
				%s
				%s
				%s
				`, color.Apply("test", color.Red), color.Apply("test2", color.Red), color.Apply("test3", color.Red)),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := SingleColoredPrinter{Color: tt.color}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
