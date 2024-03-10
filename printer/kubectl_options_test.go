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
				[37mThe following options can be passed to any command:[0m

				    [33m--add-dir-header[0m=[31mfalse[0m:
					[37mIf true, adds the file directory to the header of the log messages[0m
				    [33m--alsologtostderr[0m=[32mtrue[0m:
					[37mlog to standard error as well as files[0m
				    [33m--as[0m=[37m''[0m:
					[37mUsername to impersonate for the operation[0m
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
