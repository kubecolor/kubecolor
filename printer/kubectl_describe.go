package printer

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/kubecolor/kubecolor/config"
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
			keyColor := getColorByKeyIndent(line.KeyIndent(), basicIndentWidth, dp.TablePrinter.Theme.Describe.Key)
			key := string(line.Key)
			if withoutColon, ok := strings.CutSuffix(key, ":"); ok {
				fmt.Fprint(w, keyColor.Render(withoutColon), ":")
			} else {
				fmt.Fprint(w, keyColor.Render(key))
			}
		}
		fmt.Fprintf(w, "%s", line.Spacing)
		if len(line.Value) > 0 {
			val := string(line.Value)
			if k, v, ok := strings.Cut(val, ": "); ok { // split annotation and env from
				vColor := dp.colorize(scanner.Path(), v)
				fmt.Fprint(w, k, ": ", vColor.Render(v))
			} else if k, v, ok := strings.Cut(val, "="); ok { // split label
				vColor := dp.colorize(scanner.Path(), v)
				fmt.Fprint(w, k, "=", vColor.Render(v))
			} else {
				valColor := dp.colorize(scanner.Path(), val)
				fmt.Fprint(w, valColor.Render(val))
			}
		}
		fmt.Fprintf(w, "%s\n", line.Trailing)
	}

	if dp.tableBytes != nil {
		dp.TablePrinter.Print(dp.tableBytes, w)
		dp.tableBytes = nil
	}
}

func (dp *DescribePrinter) colorize(path describe.Path, value string) config.Color {
	value = strings.TrimSpace(value)
	pathStr := path.String()
	if col, ok := dp.colorizeStatus(value, pathStr); ok {
		return col
	}
	if col, ok := dp.colorizeArgs(value, pathStr); ok {
		return col
	}
	return getColorByValueType(value, dp.TablePrinter.Theme)
}

var describePathsToColor = []*regexp.Regexp{
	regexp.MustCompile(`^Status$`),
	regexp.MustCompile(`^(Init )?Containers/[^/]*/State(/Reason)?$`),
	regexp.MustCompile(`^Containers/[^/]*/Last State(/Reason)?$`),
}

func (dp *DescribePrinter) colorizeStatus(value, pathStr string) (config.Color, bool) {
	if !matchesAnyRegex(pathStr, describePathsToColor) {
		return config.Color{}, false
	}
	col, ok := ColorStatus(value, dp.TablePrinter.Theme)
	if !ok {
		return config.Color{}, false
	}
	return col, true
}

var describePathsToArgs = []*regexp.Regexp{
	regexp.MustCompile(`^(Init )?Containers/[^/]*/Args$`),
}

func (dp *DescribePrinter) colorizeArgs(value, pathStr string) (config.Color, bool) {
	if !matchesAnyRegex(pathStr, describePathsToArgs) {
		return config.Color{}, false
	}
	if !strings.HasPrefix(value, "-") {
		return config.Color{}, false
	}
	// Intentionally empty color, so it gets the same color as keys
	// where "--my-flag=123" args has no coloring on "--my-flag=",
	// so "--my-bool-flag" should also not be colored.
	return config.Color{}, true
}

func matchesAnyRegex(s string, regexes []*regexp.Regexp) bool {
	for _, r := range regexes {
		if r.MatchString(s) {
			return true
		}
	}
	return false
}
