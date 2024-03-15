package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_OptionsPrinter_Print(t *testing.T) {
	tests := []struct {
		name     string
		theme    *config.Theme
		input    string
		expected string
	}{
		{
			name:  "successful",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				The following options can be passed to any command:

				    --add-dir-header=false:
					If true, adds the file directory to the header of the log messages
				    --alsologtostderr=true:
					log to standard error as well as files
				    --as='':
					Username to impersonate for the operation
				`),
			expected: testutil.NewHereDoc(`
				\e[37mThe following options can be passed to any command:\e[0m

				    \e[33m--add-dir-header\e[0m=\e[31mfalse\e[0m:
					\e[37mIf true, adds the file directory to the header of the log messages\e[0m
				    \e[33m--alsologtostderr\e[0m=\e[32mtrue\e[0m:
					\e[37mlog to standard error as well as files\e[0m
				    \e[33m--as\e[0m=\e[37m''\e[0m:
					\e[37mUsername to impersonate for the operation\e[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := OptionsPrinter{Theme: tt.theme}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
