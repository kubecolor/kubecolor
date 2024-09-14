package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/internal/stringutil"
	"github.com/kubecolor/kubecolor/testutil"
)

func FuzzYAMLPrinter(f *testing.F) {
	f.Add("key: value\n")
	f.Add("  foo = \"bar\"")
	f.Add("key: |\n  foo = \"bar\"")
	f.Fuzz(func(t *testing.T, input string) {
		// treat CR as LF
		input = strings.ReplaceAll(input, "\r", "\n")
		// ignore if string is >2kB
		input = stringutil.Truncate(input, 2048)
		// make sure we have trailing newline
		if !strings.HasSuffix(input, "\n") {
			input += "\n"
		}

		printer := YAMLPrinter{
			Theme: testconfig.NullTheme,
		}
		var buf bytes.Buffer
		printer.Print(strings.NewReader(input), &buf)

		testutil.MustEqual(t, input, buf.String())
	})
}
