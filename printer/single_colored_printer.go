package printer

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"

	"github.com/kubecolor/kubecolor/config/color"
	"github.com/kubecolor/kubecolor/internal/bytesutil"
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
	scanner.Buffer(nil, bytesutil.MaxLineLength)
	for scanner.Scan() {
		fmt.Fprintf(w, "%s\n", p.Color.Render(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		slog.Error("Failed to print single-colored output.", "error", err)
	}
}
