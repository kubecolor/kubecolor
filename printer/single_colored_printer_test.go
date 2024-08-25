package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config/color"
)

func Test_SingleColoredPrinter_Print(t *testing.T) {
	input := "hello\nworld\nfoo\nbar"
	var w bytes.Buffer
	printer := SingleColoredPrinter{Color: color.MustParseColor("yellow")}
	printer.Print(strings.NewReader(input), &w)
	got := w.String()
	if got == input {
		t.Fatalf("input equals output, but colors should have been applied")
	}
}
