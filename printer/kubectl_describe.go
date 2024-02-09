package printer

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/kubecolor/kubecolor/color"
	"github.com/kubecolor/kubecolor/internal/bytesutil"
	"github.com/kubecolor/kubecolor/scanner/describe"
)

// DescribePrinter is a specific printer to print kubectl describe format.
type DescribePrinter struct {
	TablePrinter *TablePrinter

	tableBytes *bytes.Buffer
}

var onlyValuePathToColor = regexp.MustCompile(`^(Labels|Annotations|(Init )?Containers/[^/]+/Environment Variables from)/.+`)

func (dp *DescribePrinter) Print(r io.Reader, w io.Writer) {
	scanner := describe.NewScanner(r)
	const basicIndentWidth = 2 // according to kubectl describe format
	for scanner.Scan() {
		line := scanner.Line()

		if onlyValuePathToColor.MatchString(scanner.Path().String()) { // if line path matches label or annotation or env from
			line.Value = bytes.Join([][]byte{line.Key, line.Value}, line.Spacing)
			line.Key = nil
			line.Spacing = nil
		} else if bytesutil.CountColumns(line.Value, " \t") >= 3 { // when there are multiple columns, treat is as table format
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
			keyColor := getColorByKeyIndent(line.KeyIndent(), basicIndentWidth, dp.TablePrinter.Theme)
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
			if k, v, ok := strings.Cut(val, ": "); ok { // split annotation and env from
				vColor := dp.valueColor(scanner.Path(), v)
				fmt.Fprint(w, k, ": ", color.Apply(v, vColor))
			} else if k, v, ok := strings.Cut(val, "="); ok { // split label
				vColor := dp.valueColor(scanner.Path(), v)
				fmt.Fprint(w, k, "=", color.Apply(v, vColor))
			} else {
				valColor := dp.valueColor(scanner.Path(), val)
				fmt.Fprint(w, color.Apply(val, valColor))
			}
		}
		fmt.Fprintf(w, "%s\n", line.Trailing)
	}

	if dp.tableBytes != nil {
		dp.TablePrinter.Print(dp.tableBytes, w)
		dp.tableBytes = nil
	}
}

func (dp *DescribePrinter) valueColor(path describe.Path, value string) color.Color {
	value = strings.TrimSpace(value)
	if describeUseStatusColoring(path) {
		if col, ok := ColorStatus(value, dp.TablePrinter.Theme); ok {
			return col
		}
	}
	return getColorByValueType(value, dp.TablePrinter.Theme)
}

var describePathsToColor = []*regexp.Regexp{
	regexp.MustCompile(`^Status$`),
	regexp.MustCompile(`^(Init )?Containers/[^/]*/State(/Reason)?$`),
	regexp.MustCompile(`^Containers/[^/]*/Last State(/Reason)?$`),
}

func describeUseStatusColoring(path describe.Path) bool {
	if len(path) == 0 {
		return false
	}
	str := path.String()

	for _, r := range describePathsToColor {
		if r.MatchString(str) {
			return true
		}
	}
	return false
}
