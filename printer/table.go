package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/scanner/tablescan"
)

// TablePrinter is used in table output such as "kubectl get"
// and the events tables in "kubectl describe"
type TablePrinter struct {
	WithHeader     bool
	DarkBackground bool
	Theme          *config.Theme
	ColumnFilter   func(columnIndex int, column string) string

	hasLeadingNamespaceColumn bool
}

// ensures it implements the interface
var _ Printer = &TablePrinter{}

func NewTablePrinter(withHeader bool, theme *config.Theme, columnFilter func(columnIndex int, column string) string) *TablePrinter {
	return &TablePrinter{
		WithHeader:   withHeader,
		Theme:        theme,
		ColumnFilter: columnFilter,
	}
}

// Print implements [Printer.Print]
func (p *TablePrinter) Print(r io.Reader, w io.Writer) {
	isFirstLine := true
	scanner := tablescan.NewScanner(r)
	for scanner.Scan() {
		cells := scanner.Cells()
		if len(cells) == 0 {
			fmt.Fprint(w, "\n")
			continue
		}
		peekNextLine, hasNextLine := scanner.PeekText()
		if (p.WithHeader && isFirstLine) ||
			isAllUpper(scanner.Text()) ||
			(hasNextLine && isOnlySymbols(peekNextLine)) ||
			isOnlySymbols(scanner.Text()) {

			isFirstLine = false
			leadingSpaces := scanner.LeadingSpaces()
			withoutSpaces := scanner.Text()[len(leadingSpaces):]
			fmt.Fprintf(w, "%s%s\n", leadingSpaces, p.Theme.Table.Header.Render(withoutSpaces))

			if strings.EqualFold(cells[0].Trimmed, "namespace") {
				p.hasLeadingNamespaceColumn = true
			}
			continue
		}

		fmt.Fprintf(w, "%s", scanner.LeadingSpaces())
		p.printLineAsTableFormat(w, cells, p.Theme.Table.Columns)
	}
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
func (p *TablePrinter) printLineAsTableFormat(w io.Writer, cells []tablescan.Cell, colorsPreset []config.Color) {
	for i, cell := range cells {
		c := p.getColumnBaseColor(i, colorsPreset)

		cellText := cell.Trimmed
		if p.ColumnFilter != nil {
			cellText = p.ColumnFilter(i, cellText)
		}
		// Write colored column
		if cellText != "" {
			fmt.Fprint(w, c.Render(cellText))
		}
		fmt.Fprint(w, cell.TrailingSpaces)
	}

	fmt.Fprintf(w, "\n")
}

func (p *TablePrinter) getColumnBaseColor(index int, colorsPreset []config.Color) config.Color {
	if len(colorsPreset) == 0 {
		return config.Color{}
	}
	if p.hasLeadingNamespaceColumn {
		index--
		if index < 0 {
			index += len(colorsPreset)
		}
	}
	return colorsPreset[index%len(colorsPreset)]
}
