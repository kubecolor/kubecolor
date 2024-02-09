package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/color"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_ApplyPrinter_Print(t *testing.T) {
	tests := []struct {
		name        string
		themePreset color.Preset
		input       string
		expected    string
	}{
		{
			name:        "created",
			themePreset: color.PresetDark,
			input: testutil.NewHereDoc(`
				deployment.apps/foo created`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo [32mcreated[0m
			`),
		},
		{
			name:        "configured",
			themePreset: color.PresetDark,
			input: testutil.NewHereDoc(`
				deployment.apps/foo configured`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo [33mconfigured[0m
			`),
		},
		{
			name:        "unchanged",
			themePreset: color.PresetDark,
			input: testutil.NewHereDoc(`
				deployment.apps/foo unchanged`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo [35munchanged[0m
			`),
		},
		{
			name:        "client dry run",
			themePreset: color.PresetDark,
			input: testutil.NewHereDoc(`
				deployment.apps/foo unchanged (dry run)`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo [35munchanged[0m [36m(dry run)[0m
			`),
		},
		{
			name:        "server dry run",
			themePreset: color.PresetDark,
			input: testutil.NewHereDoc(`
				deployment.apps/foo unchanged (server dry run)`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo [35munchanged[0m [36m(server dry run)[0m
			`),
		},
		{
			name:        "something else. This likely won't happen but fallbacks here just in case.",
			themePreset: color.PresetDark,
			input: testutil.NewHereDoc(`
				deployment.apps/foo bar`),
			expected: testutil.NewHereDoc(`
				[32mdeployment.apps/foo bar[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := ApplyPrinter{Theme: color.NewTheme(tt.themePreset)}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
