package printer

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/color"
	"github.com/kubecolor/kubecolor/scanner/describe"
)

// ExplainPrinter is a specific printer to print kubectl explain format.
type ExplainPrinter struct {
	DarkBackground bool
	ColorSchema    ColorSchema
	Recursive      bool
}

func (ep *ExplainPrinter) Print(r io.Reader, w io.Writer) {
	scanner := describe.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Line()
		isFields := scanner.Path().HasPrefix("FIELDS")

		fmt.Fprintf(w, "%s", line.Indent)
		if bytes.ContainsAny(line.Key, " \t-.") {
			fmt.Fprintf(w, "%s", line.Key)
		} else if len(line.Key) > 0 {
			keyColor := ep.keyColor(line, isFields)
			key := string(line.Key)
			if withoutColon, ok := strings.CutSuffix(key, ":"); ok {
				fmt.Fprint(w, color.Apply(withoutColon, keyColor), ":")
			} else {
				fmt.Fprint(w, color.Apply(key, keyColor))
			}
		}
		fmt.Fprintf(w, "%s", line.Spacing)
		ep.printVal(w, string(line.Value))
		fmt.Fprintf(w, "%s\n", line.Trailing)
	}
}

func (ep *ExplainPrinter) keyColor(line describe.Line, isFields bool) color.Color {
	if ep.Recursive && isFields {
		return getColorByKeyIndent(line.KeyIndent(), 2, ep.ColorSchema)
	}

	return getColorByKeyIndent(0, 2, ep.ColorSchema)
}

func (ep *ExplainPrinter) printVal(w io.Writer, val string) {
	const suffix = "-required-"
	if withoutSuffix, ok := strings.CutSuffix(val, suffix); ok {
		fmt.Fprint(w, withoutSuffix)
		fmt.Fprint(w, color.Apply(suffix, ep.ColorSchema.RequiredColor))
		return
	}
	fmt.Fprint(w, val)
}
