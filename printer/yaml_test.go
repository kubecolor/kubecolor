package printer

import (
	"bytes"
	"errors"
	"testing"

	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

func TestYAMLPrinter_fail(t *testing.T) {
	var logBuf bytes.Buffer
	testutil.SetTestLogger(t, &logBuf)

	var outBuf bytes.Buffer
	printer := YAMLPrinter{Theme: testconfig.NullTheme}
	printer.Print(testutil.DummyReader{ReadFunc: func(b []byte) (int, error) { return 0, errors.New("test") }}, &outBuf)

	testutil.Equal(t, "", outBuf.String(), "output")
	testutil.Equal(t, "level=ERROR msg=\"Failed to print YAML output.\" error=test\n", logBuf.String(), "logs")
}

func TestYAMLPrinter_colorizeBlockScalarHeader(t *testing.T) {
	tests := []struct {
		input      string
		wantOK     bool
		wantIndent int
	}{
		{input: "|", wantOK: true},
		{input: ">-", wantOK: true},
		{input: "|+", wantOK: true},
		{input: "|2", wantOK: true, wantIndent: 2},
		{input: ">+4", wantOK: true, wantIndent: 4},
		{input: "|2-", wantOK: true, wantIndent: 2},
		{input: "", wantOK: false},
		{input: "foo", wantOK: false},
		{input: "|0", wantOK: false},
		{input: "|2x", wantOK: false},
		{input: "|22", wantOK: false},
		{input: "|+-", wantOK: false},
		{input: "|-2-", wantOK: false},
		{input: "| foo", wantOK: false},
	}

	printer := YAMLPrinter{Theme: testconfig.NullTheme}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			header, indent, ok := printer.colorizeBlockScalarHeader(tc.input)
			testutil.Equal(t, tc.wantOK, ok, "ok")
			testutil.Equal(t, tc.wantIndent, indent, "indent")
			if tc.wantOK {
				testutil.Equal(t, tc.input, header, "header")
			} else {
				testutil.Equal(t, "", header, "header")
			}
		})
	}
}
