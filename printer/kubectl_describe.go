package printer

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"

	"github.com/kubecolor/kubecolor/internal/bytesutil"
	"github.com/kubecolor/kubecolor/scanner/describe"
)

// DescribePrinter is used on "kubectl describe" output
type DescribePrinter struct {
	TablePrinter *TablePrinter

	tableBytes *bytes.Buffer
}

// ensures it implements the interface
var _ Printer = &DescribePrinter{}

var onlyValuePathToColor = regexp.MustCompile(`^(Labels|Annotations|(Init )?Containers/[^/]+/Environment Variables from)/.+`)

// Print implements [Printer.Print]
func (p *DescribePrinter) Print(r io.Reader, w io.Writer) {
	scanner := describe.NewScanner(r)
	const basicIndentWidth = 2 // according to kubectl describe format
	for scanner.Scan() {
		line := scanner.Line()

		if onlyValuePathToColor.MatchString(scanner.Path().String()) { // if line path matches label or annotation or env from
			line.Value = bytes.Join([][]byte{line.Key, line.Value}, line.Spacing)
			line.Key = nil
			line.Spacing = nil
		} else if bytesutil.CountColumns(line.Value, " \t") >= 3 { // when there are multiple columns, treat is as table format
			if p.tableBytes == nil {
				p.tableBytes = &bytes.Buffer{}
			}
			fmt.Fprintln(p.tableBytes, line.String())
			continue
		} else if p.tableBytes != nil {
			p.TablePrinter.Print(p.tableBytes, w)
			p.tableBytes = nil
		}

		fmt.Fprintf(w, "%s", line.Indent)
		if len(line.Key) > 0 {
			keyColor := ColorDataKey(line.KeyIndent(), basicIndentWidth, p.TablePrinter.Theme.Describe.Key)
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
				fmt.Fprint(w, k, ": ", p.colorize(scanner.Path(), v))
			} else if k, v, ok := strings.Cut(val, "="); ok { // split label
				fmt.Fprint(w, k, "=", p.colorize(scanner.Path(), v))
			} else {
				fmt.Fprint(w, p.colorize(scanner.Path(), val))
			}
		}
		fmt.Fprintf(w, "%s\n", line.Trailing)
	}
	if err := scanner.Err(); err != nil {
		slog.Error("Failed to print describe output.", "error", err)
	}

	if p.tableBytes != nil {
		p.TablePrinter.Print(p.tableBytes, w)
		p.tableBytes = nil
	}
}

func (p *DescribePrinter) colorize(path describe.Path, value string) string {
	value = strings.TrimSpace(value)
	pathStr := path.String()
	if col, ok := p.colorizeStatus(value, pathStr); ok {
		return col
	}
	if col, ok := p.colorizeArgs(value, pathStr); ok {
		return col
	}
	return ColorDataValue(value, p.TablePrinter.Theme).Render(value)
}

var describePathsToColor = []*regexp.Regexp{
	regexp.MustCompile(`^Status$`),
	regexp.MustCompile(`^(Init )?Containers/[^/]*/State(/Reason)?$`),
	regexp.MustCompile(`^Containers/[^/]*/Last State(/Reason)?$`),
}

func (p *DescribePrinter) colorizeStatus(value, pathStr string) (string, bool) {
	if !matchesAnyRegex(pathStr, describePathsToColor) {
		return value, false
	}
	col, ok := ColorStatus(value, p.TablePrinter.Theme)
	if !ok {
		return value, false
	}
	return col, true
}

var describePathsToArgs = []*regexp.Regexp{
	regexp.MustCompile(`^(Init )?Containers/[^/]*/Args$`),
}

func (p *DescribePrinter) colorizeArgs(value, pathStr string) (string, bool) {
	if !matchesAnyRegex(pathStr, describePathsToArgs) {
		return value, false
	}
	if !strings.HasPrefix(value, "-") {
		return value, false
	}
	// Intentionally empty color, so it gets the same color as keys
	// where "--my-flag=123" args has no coloring on "--my-flag=",
	// so "--my-bool-flag" should also not be colored.
	return value, true
}

func matchesAnyRegex(s string, regexes []*regexp.Regexp) bool {
	for _, r := range regexes {
		if r.MatchString(s) {
			return true
		}
	}
	return false
}
