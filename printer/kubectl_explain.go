package printer

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/scanner/describe"
)

// ExplainPrinter is used in "kubectl explain" output
type ExplainPrinter struct {
	Theme     *config.Theme
	Recursive bool
}

// ensures it implements the interface
var _ Printer = &ExplainPrinter{}

// Print implements [Printer.Print]
func (p *ExplainPrinter) Print(r io.Reader, w io.Writer) {
	scanner := describe.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Line()
		isFields := scanner.Path().HasPrefix("FIELDS")

		fmt.Fprintf(w, "%s", line.Indent)
		if bytes.ContainsAny(line.Key, " \t-.") {
			fmt.Fprintf(w, "%s", line.Key)
		} else if len(line.Key) > 0 {
			keyColor := p.keyColor(line, isFields)
			key := string(line.Key)
			if withoutColon, ok := strings.CutSuffix(key, ":"); ok {
				fmt.Fprint(w, keyColor.Render(withoutColon), ":")
			} else {
				fmt.Fprint(w, keyColor.Render(key))
			}
		}
		fmt.Fprintf(w, "%s", line.Spacing)
		p.printVal(w, string(line.Value))
		fmt.Fprintf(w, "%s\n", line.Trailing)
	}
}

func (p *ExplainPrinter) keyColor(line describe.Line, isFields bool) config.Color {
	if p.Recursive && isFields {
		return ColorDataKey(line.KeyIndent(), 2, p.Theme.Explain.Key)
	}

	return ColorDataKey(0, 2, p.Theme.Explain.Key)
}

func (p *ExplainPrinter) printVal(w io.Writer, val string) {
	const suffix = "-required-"
	if withoutSuffix, ok := strings.CutSuffix(val, suffix); ok {
		fmt.Fprint(w, withoutSuffix)
		fmt.Fprint(w, p.Theme.Explain.Required.Render(suffix))
		return
	}
	fmt.Fprint(w, val)
}
