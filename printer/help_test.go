package printer

import (
	"bytes"
	"errors"
	"testing"

	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

func TestHelpPrinter_fail(t *testing.T) {
	var logBuf bytes.Buffer
	testutil.SetTestLogger(t, &logBuf)

	var outBuf bytes.Buffer
	printer := HelpPrinter{Theme: testconfig.NullTheme}
	printer.Print(testutil.DummyReader{ReadFunc: func(b []byte) (int, error) { return 0, errors.New("test") }}, &outBuf)

	testutil.Equal(t, "", outBuf.String(), "output")
	testutil.Equal(t, "level=ERROR msg=\"Failed to print help output.\" error=test\n", logBuf.String(), "logs")
}
