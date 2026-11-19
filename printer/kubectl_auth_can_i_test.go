package printer

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config/color"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_CanIPrinter_Print_fallback(t *testing.T) {
	input := "Warning: resource 'foos' is not namespaced"
	var w bytes.Buffer
	printer := CanIPrinter{
		AllowedColor:  color.MustParse("green"),
		DeniedColor:   color.MustParse("red"),
		FallbackColor: color.MustParse("yellow"),
	}
	printer.Print(strings.NewReader(input), &w)
	got := w.String()
	if got == input {
		t.Fatalf("input equals output, but colors should have been applied")
	}
}

func TestCanIPrinter_fail(t *testing.T) {
	var logBuf bytes.Buffer
	testutil.SetTestLogger(t, &logBuf)

	var outBuf bytes.Buffer
	printer := CanIPrinter{}
	printer.Print(testutil.DummyReader{ReadFunc: func(b []byte) (int, error) { return 0, errors.New("test") }}, &outBuf)

	testutil.Equal(t, "", outBuf.String(), "output")
	testutil.Equal(t, "level=ERROR msg=\"Failed to print can-i output.\" error=test\n", logBuf.String(), "logs")
}
