package printer

import (
	"bytes"
	"errors"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

func TestLogsPrinter_fail(t *testing.T) {
	var logBuf bytes.Buffer
	testutil.SetTestLogger(t, &logBuf)

	var outBuf bytes.Buffer
	printer := LogsPrinter{}
	printer.Print(testutil.DummyReader{ReadFunc: func(b []byte) (int, error) { return 0, errors.New("test") }}, &outBuf)

	testutil.Equal(t, "", outBuf.String(), "output")
	testutil.Equal(t, "level=ERROR msg=\"Failed to print log output.\" error=test\n", logBuf.String(), "logs")
}
