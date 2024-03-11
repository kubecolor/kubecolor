package printer

import (
	"fmt"
	"io"
	"slices"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/scanner/describe"
)

type OptionsPrinter struct {
	Theme *config.Theme
}

func (op *OptionsPrinter) Print(r io.Reader, w io.Writer) {
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
				op.Theme.Options.Flag.Render(string(line.Key)),
				line.Spacing,
				getColorByValueType(val, op.Theme).Render(val),
				line.Trailing)
			continue
		}

		fmt.Fprintf(w, "%s%s%s\n",
			line.Indent,
			op.Theme.Data.String.Render(string(slices.Concat(line.Key, line.Spacing, line.Value))),
			line.Trailing)
	}
}
