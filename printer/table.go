package printer

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/scanner/tablescan"
)

type TablePrinter struct {
	WithHeader     bool
	DarkBackground bool
	Theme          *config.Theme
	ColorDeciderFn func(index int, column string) (config.Color, bool)

	hasLeadingNamespaceColumn bool
}

func NewTablePrinter(withHeader bool, theme *config.Theme, colorDeciderFn func(index int, column string) (config.Color, bool)) *TablePrinter {
	return &TablePrinter{
		WithHeader:     withHeader,
		Theme:          theme,
		ColorDeciderFn: colorDeciderFn,
	}
}

func (tp *TablePrinter) Print(r io.Reader, w io.Writer) {
	isFirstLine := true
	scanner := tablescan.NewScanner(r)
	for scanner.Scan() {
		cells := scanner.Cells()
		if len(cells) == 0 {
			fmt.Fprint(w, "\n")
			continue
		}
		if (tp.WithHeader && isFirstLine) || isAllCellsUpper(cells) {
			isFirstLine = false
			fmt.Fprintf(w, "%s\n", tp.Theme.Header.Render(scanner.Text()))

			if strings.EqualFold(cells[0].Trimmed, "namespace") {
				tp.hasLeadingNamespaceColumn = true
			}
			continue
		}

		fmt.Fprintf(w, "%s", scanner.LeadingSpaces())
		tp.printLineAsTableFormat(w, cells, tp.Theme.ColumnCycle)
	}
}

func isAllCellsUpper(cells []tablescan.Cell) bool {
	for _, c := range cells {
		if !isAllUpper(c.Trimmed) {
			return false
		}
	}
	return true
}

func isAllUpper(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

// printTableFormat prints a line to w in kubectl "table" Format.
// Table format is something like:
//
//	NAME                     READY   STATUS    RESTARTS   AGE
//	nginx-6799fc88d8-dnmv5   1/1     Running   0          31h
//	nginx-6799fc88d8-m8pbc   1/1     Running   0          31h
//	nginx-6799fc88d8-qdf9b   1/1     Running   0          31h
//	nginx-8spn9              1/1     Running   0          31h
//	nginx-dplns              1/1     Running   0          31h
//	nginx-lpv5x              1/1     Running   0          31h
func (tp *TablePrinter) printLineAsTableFormat(w io.Writer, cells []tablescan.Cell, colorsPreset []config.Color) {
	for i, cell := range cells {
		c := tp.getColumnBaseColor(i, colorsPreset)

		if tp.ColorDeciderFn != nil {
			if cc, ok := tp.ColorDeciderFn(i, cell.Trimmed); ok {
				c = cc // prior injected deciderFn result
			}
		}
		// Write colored column
		if cell.Trimmed != "" {
			fmt.Fprint(w, c.Render(cell.Trimmed))
		}
		fmt.Fprint(w, cell.TrailingSpaces)
	}

	fmt.Fprintf(w, "\n")
}

func (tp *TablePrinter) getColumnBaseColor(index int, colorsPreset []config.Color) config.Color {
	if tp.hasLeadingNamespaceColumn {
		index++
	}
	return colorsPreset[index%len(colorsPreset)]
}
