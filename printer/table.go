package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/color"
)

type TablePrinter struct {
	WithHeader     bool
	DarkBackground bool
	ColorDeciderFn func(index int, column string) (color.Color, bool)

	isFirstLine   bool
	indexColorMap map[int]color.Color
	tempColors    []color.Color
}

func NewTablePrinter(withHeader, darkBackground bool, colorDeciderFn func(index int, column string) (color.Color, bool)) *TablePrinter {
	return &TablePrinter{
		WithHeader:     withHeader,
		DarkBackground: darkBackground,
		ColorDeciderFn: colorDeciderFn,
		indexColorMap:  map[int]color.Color{},
		tempColors:     []color.Color{},
	}
}

func (tp *TablePrinter) Print(r io.Reader, w io.Writer) {
	tp.isFirstLine = true
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if tp.isHeader(line) {
			fmt.Fprintf(w, "%s\n", color.Apply(line, getHeaderColorByBackground(tp.DarkBackground)))
			tp.isFirstLine = false
			continue
		}

		tp.printLineAsTableFormat(w, line, getColorsByBackground(tp.DarkBackground))
	}
}

func (tp *TablePrinter) isHeader(line string) bool {
	// If every character is upper case, probably it's a header line.
	// e.g.
	// kubecolor get pod,rs
	// NAME                         READY   STATUS    RESTARTS   AGE
	// pod/nginx-8spn9              1/1     Running   1          19d
	// pod/nginx-dplns              1/1     Running   1          19d
	// pod/nginx-lpv5x              1/1     Running   1          19d

	// NAME                               DESIRED   CURRENT   READY   AGE <- this
	// replicaset.apps/nginx              3         3         3       19d
	// replicaset.apps/nginx-6799fc88d8   3         3         3       19d
	isEveryCharacterUpper := strings.ToUpper(line) == line
	return (tp.WithHeader && tp.isFirstLine) || isEveryCharacterUpper
}

// printTableFormat prints a line to w in kubectl "table" Format.
// Table format is something like:
// ---------------------------------------------------------
// NAME                     READY   STATUS    RESTARTS   AGE
// nginx-6799fc88d8-dnmv5   1/1     Running   0          31h
// nginx-6799fc88d8-m8pbc   1/1     Running   0          31h
// nginx-6799fc88d8-qdf9b   1/1     Running   0          31h
// nginx-8spn9              1/1     Running   0          31h
// nginx-dplns              1/1     Running   0          31h
// nginx-lpv5x              1/1     Running   0          31h
// ---------------------------------------------------------
// This function requires a line and tries to colorize it by each column.
// If dark is true, use colors which are readable in dark-backgrounded environment, else,
// it uses colors for light-backgrounded env.
// This function doesn't respect if the line is "header", so
// if you want to specify a special color for header, you must not pass the line
// to this function.
// deciderFn is a function to return context-specific color to be used to decorate a column.
// If the function returned ok=true, then returned color will be used for the column.
// If it returned ok=false, then default configurated color will be used.
// If deciderFn is null, then this function uses the default configurated color.
func (tp *TablePrinter) printLineAsTableFormat(w io.Writer, line string, colorsPreset []color.Color) {
	columns := spaces.Split(line, -1)
	spacesIndices := spaces.FindAllStringIndex(line, -1)

	if len(columns) == len(spacesIndices)-1 {
		// It should not come here.
		panic("kubecolor: unexpected format as table. this must be a bug of kubecolor")
	}

	for i, column := range columns {
		index := 0
		if i != 0 {
			index = spacesIndices[i-1][1] + 1
		}

		c := tp.decideColorForTable(index, colorsPreset)
		if tp.ColorDeciderFn != nil {
			if cc, ok := tp.ColorDeciderFn(i, column); ok {
				c = cc // prior injected deciderFn result
			}
		}
		// Write colored column
		fmt.Fprintf(w, "%s", color.Apply(column, c))
		// Write spaces based on actual output
		// When writing the most left column, no extra spaces needed.
		if i <= len(spacesIndices)-1 {
			spacesIndex := spacesIndices[i]
			fmt.Fprintf(w, "%s", toSpaces(spacesIndex[1]-spacesIndex[0]))
		}
	}

	fmt.Fprintf(w, "\n")
}

func (tp *TablePrinter) decideColorForTable(index int, colors []color.Color) color.Color {
	if len(tp.tempColors) == 0 {
		tp.tempColors = make([]color.Color, len(colors))
		copy(tp.tempColors, colors)
	}

	if c, ok := tp.indexColorMap[index]; ok {
		return c
	}

	c := tp.tempColors[0]
	tp.indexColorMap[index] = c
	tp.tempColors = tp.tempColors[1:]

	return c
}
