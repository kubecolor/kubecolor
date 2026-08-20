package printer

import (
	"bytes"
	"errors"
	"testing"

	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

func TestHelpPrinter_colorizeUrls(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "bracketed url keeps last character",
			input: "See [https://kubernetes.io/docs/reference] for more.",
			want:  "See [https://kubernetes.io/docs/reference] for more.",
		},
		{
			name:  "bare url",
			input: "See https://kubernetes.io/docs/reference for more.",
			want:  "See https://kubernetes.io/docs/reference for more.",
		},
		{
			name:  "multiple bracketed urls",
			input: "a [http://golang.org/pkg/text/template/#pkg-overview] b [https://kubernetes.io/docs/reference/kubectl/jsonpath/] c",
			want:  "a [http://golang.org/pkg/text/template/#pkg-overview] b [https://kubernetes.io/docs/reference/kubectl/jsonpath/] c",
		},
	}

	printer := HelpPrinter{Theme: testconfig.NullTheme}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testutil.Equal(t, tc.want, printer.colorizeUrls(tc.input), "colorized urls")
		})
	}
}

func TestHelpPrinter_fail(t *testing.T) {
	var logBuf bytes.Buffer
	testutil.SetTestLogger(t, &logBuf)

	var outBuf bytes.Buffer
	printer := HelpPrinter{Theme: testconfig.NullTheme}
	printer.Print(testutil.DummyReader{ReadFunc: func(b []byte) (int, error) { return 0, errors.New("test") }}, &outBuf)

	testutil.Equal(t, "", outBuf.String(), "output")
	testutil.Equal(t, "level=ERROR msg=\"Failed to print help output.\" error=test\n", logBuf.String(), "logs")
}
