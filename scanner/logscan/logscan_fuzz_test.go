package logscan

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/internal/stringutil"
	"github.com/kubecolor/kubecolor/testutil"
)

func FuzzScanner(f *testing.F) {
	f.Add("INFO: hello world\n")
	f.Fuzz(func(t *testing.T, input string) {
		// treat CR as LF
		input = strings.ReplaceAll(input, "\r", "\n")
		// ignore if string is >2kB
		input = stringutil.Truncate(input, 2048)
		// make sure we have trailing newline
		if !strings.HasSuffix(input, "\n") {
			input += "\n"
		}

		scanner := NewScanner(strings.NewReader(input))

		var buf bytes.Buffer
		for scanner.Scan() {
			buf.WriteString(scanner.Token().Text)
		}

		testutil.MustEqual(t, input, buf.String())
	})
}
