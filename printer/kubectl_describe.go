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
				vColor := p.colorize(scanner.Path(), v)
				fmt.Fprint(w, k, ": ", vColor.Render(v))
			} else if k, v, ok := strings.Cut(val, "="); ok { // split label
				vColor := p.colorize(scanner.Path(), v)
				fmt.Fprint(w, k, "=", vColor.Render(v))
			} else {
				valColor := p.colorize(scanner.Path(), val)
				fmt.Fprint(w, valColor.Render(val))
			}
		}
		fmt.Fprintf(w, "%s\n", line.Trailing)
	}

	if p.tableBytes != nil {
		p.TablePrinter.Print(p.tableBytes, w)
		p.tableBytes = nil
	}
}

func (p *DescribePrinter) colorize(path describe.Path, value string) config.Color {
	value = strings.TrimSpace(value)
	pathStr := path.String()
	if col, ok := p.colorizeStatus(value, pathStr); ok {
		return col
	}
	if col, ok := p.colorizeArgs(value, pathStr); ok {
		return col
	}
	return ColorDataValue(value, p.TablePrinter.Theme)
}

var describePathsToColor = []*regexp.Regexp{
	regexp.MustCompile(`^Status$`),
	regexp.MustCompile(`^(Init )?Containers/[^/]*/State(/Reason)?$`),
	regexp.MustCompile(`^Containers/[^/]*/Last State(/Reason)?$`),
}

func (p *DescribePrinter) colorizeStatus(value, pathStr string) (config.Color, bool) {
	if !matchesAnyRegex(pathStr, describePathsToColor) {
		return config.Color{}, false
	}
	col, ok := ColorStatus(value, p.TablePrinter.Theme)
	if !ok {
		return config.Color{}, false
	}
	return col, true
}

var describePathsToArgs = []*regexp.Regexp{
	regexp.MustCompile(`^(Init )?Containers/[^/]*/Args$`),
}

func (p *DescribePrinter) colorizeArgs(value, pathStr string) (config.Color, bool) {
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
