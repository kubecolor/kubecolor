package printer

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/color"
	"github.com/kubecolor/kubecolor/scanner/describe"
)

// DescribePrinter is a specific printer to print kubectl describe format.
type DescribePrinter struct {
	DarkBackground bool
	TablePrinter   *TablePrinter

	tableBytes *bytes.Buffer
}

func (dp *DescribePrinter) Print(r io.Reader, w io.Writer) {
	scanner := describe.NewScanner(r)
	const basicIndentWidth = 2 // according to kubectl describe format
	doubleSpace := []byte("  ")
	for scanner.Scan() {
		line := scanner.Line()

		// when there are multiple columns, treat is as table format
		if bytes.Count(line.Value, doubleSpace) > 3 {
			if dp.tableBytes == nil {
				dp.tableBytes = &bytes.Buffer{}
			}
			fmt.Fprintln(dp.tableBytes, line.String())
			continue
		} else if dp.tableBytes != nil {
			dp.TablePrinter.Print(dp.tableBytes, w)
			dp.tableBytes = nil
		}

		fmt.Fprintf(w, "%s", line.Indent)
		if len(line.Key) > 0 {
			keyColor := getColorByKeyIndent(line.KeyIndent(), basicIndentWidth, dp.DarkBackground)
			key := string(line.Key)
			if withoutColon, ok := strings.CutSuffix(key, ":"); ok {
				fmt.Fprint(w, color.Apply(withoutColon, keyColor), ":")
			} else {
				fmt.Fprint(w, color.Apply(key, keyColor))
			}
		}
		fmt.Fprintf(w, "%s", line.Spacing)
		if len(line.Value) > 0 {
			val := string(line.Value)
			valColor := getColorByValueType(val, dp.DarkBackground)
			fmt.Fprint(w, color.Apply(val, valColor))
		}
		fmt.Fprintf(w, "%s\n", line.Trailing)
	}
}
