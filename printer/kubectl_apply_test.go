package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_ApplyPrinter_Print(t *testing.T) {
	tests := []struct {
		name     string
		theme    *config.Theme
		input    string
		expected string
	}{
		{
			name:  "created",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				deployment.apps/foo created`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo \e[32mcreated\e[0m
			`),
		},
		{
			name:  "configured",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				deployment.apps/foo configured`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo \e[33mconfigured\e[0m
			`),
		},
		{
			name:  "unchanged",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				deployment.apps/foo unchanged`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo \e[35munchanged\e[0m
			`),
		},
		{
			name:  "client dry run",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				deployment.apps/foo unchanged (dry run)`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo \e[35munchanged\e[0m \e[36m(dry run)\e[0m
			`),
		},
		{
			name:  "server dry run",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				deployment.apps/foo unchanged (server dry run)`),
			expected: testutil.NewHereDoc(`
				deployment.apps/foo \e[35munchanged\e[0m \e[36m(server dry run)\e[0m
			`),
		},
		{
			name:  "something else. This likely won't happen but fallbacks here just in case.",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				deployment.apps/foo bar`),
			expected: testutil.NewHereDoc(`
				\e[32mdeployment.apps/foo bar\e[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := ApplyPrinter{Theme: tt.theme}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
