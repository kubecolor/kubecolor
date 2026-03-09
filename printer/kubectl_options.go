package printer

import (
	"fmt"
	"io"
	"log/slog"
	"slices"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/scanner/describe"
)

// OptionsPrinter is used in "kubectl options" output
type OptionsPrinter struct {
	Theme *config.Theme
}

// ensures it implements the interface
var _ Printer = &OptionsPrinter{}

// Print implements [Printer.Print]
func (p *OptionsPrinter) Print(r io.Reader, w io.Writer) {
	scanner := describe.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Line()

		if line.IsZero() {
			fmt.Fprintln(w)
			continue
		}

		if len(scanner.Path()) == 2 {
			val := string(line.Value)
			fmt.Fprintf(w, "%s%s%s%s%s\n",
				line.Indent,
				p.Theme.Options.Flag.Render(string(line.Key)),
				line.Spacing,
				ColorDataValue(val, p.Theme).Render(val),
				line.Trailing)
			continue
		}

		fmt.Fprintf(w, "%s%s%s\n",
			line.Indent,
			p.Theme.Data.String.Render(string(slices.Concat(line.Key, line.Spacing, line.Value))),
			line.Trailing)
	}
	if err := scanner.Err(); err != nil {
		slog.Error("Failed to print options output.", "error", err)
	}
}
