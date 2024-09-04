package printer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kubecolor/kubecolor/config/color"
)

// SingleColoredPrinter is a printer to print something
// using only the single pre-configured color.
type SingleColoredPrinter struct {
	Color color.Color
}

// ensures it implements the interface
var _ Printer = &SingleColoredPrinter{}

// Print implements [Printer.Print]
func (p *SingleColoredPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Fprintf(w, "%s\n", p.Color.Render(scanner.Text()))
	}
}
