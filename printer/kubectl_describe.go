package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/color"
)

// DescribePrinter is a specific printer to print kubectl describe format.
type DescribePrinter struct {
	DarkBackground bool
	TablePrinter   *TablePrinter
}

func (dp *DescribePrinter) Print(r io.Reader, w io.Writer) {
	basicIndentWidth := 2 // according to kubectl describe format
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			fmt.Fprintln(w)
			continue
		}

		// Split a line by spaces to colorize and render them
		// For example:
		// e.g. 1-----------------
		// Status:         Running
		// -----------------------
		// spacesIndices: [[7, 15]] // <- where spaces locate
		// columns: ["Status:", "Running"]
		//
		// e.g. 2--------------------------------------------
		//     Ports:          10001/TCP, 5000/TCP, 18000/TCP
		// --------------------------------------------------
		// spacesIndices: [[0, 3], [10, 19]] // <- where spaces locate
		// columns: ["Ports:", "10001/TCP, 5000/TCP, 18000/TCP"]
		//
		// So now, we know where to render which column.
		spacesIndices := spaces.FindAllStringIndex(line, -1)
		columns := spaces.Split(line, -1)
		// when the line has indent (spaces on left), the first item will be
		// just a "" and we don't need it so remove
		if len(columns) > 0 {
			if columns[0] == "" {
				columns = columns[1:]
			}
		}

		// First, identify if there is an indent
		indentCnt := findIndent(line)
		indent := toSpaces(indentCnt)
		if indentCnt > 0 {
			// TODO: Remove this condition for workaround
			// Basically, kubectl describe output has its indentation level
			// with **2** spaces, but "Resource Quota" section in
			// `kubectl describe ns` output has only 1 space at the head.
			// Because of it, indentCnt is still 1, but the indent space is not in `spacesIndices` (see regex definition of `spaces`)
			// So it must be checked here
			// https://github.com/kubecolor/kubecolor/issues/36
			// When https://github.com/kubernetes/kubectl/issues/1005#issuecomment-758385759 is fixed
			// this is not needed anymore.
			if indentCnt > 1 {
				// when an indent exists, removes it because it's already captured by "indent" var
				spacesIndices = spacesIndices[1:]
			}
		}

		// when there are multiple columns, treat is as table format
		if len(columns) > 2 {
			dp.TablePrinter.printLineAsTableFormat(w, line, getColorsByBackground(dp.DarkBackground))
			continue
		}

		// First, write the first value assuming it's a key
		keyColor := getColorByKeyIndent(indentCnt, basicIndentWidth, dp.DarkBackground)
		valColor := getColorByValueType(columns[0], dp.DarkBackground)

		// TODO: Remove this if statement for workaround
		// Basically, kubectl describe output has its indentation level
		// with **2** spaces, but "Resource Quota" section in
		// `kubectl describe ns` output has only 1 space at the head.
		// Because of it, `spaces.Split` doesn't trim the head space (see the regex definition of `spaces`)
		// So it must be trimmed here
		// https://github.com/kubecolor/kubecolor/issues/36
		// When https://github.com/kubernetes/kubectl/issues/1005#issuecomment-758385759 is fixed
		// this is not needed anymore.
		if strings.HasPrefix(columns[0], " ") {
			columns[0] = strings.TrimLeft(columns[0], " ")
		}

		if strings.HasSuffix(columns[0], ":") {
			// trailing ":" should not be colorized
			fmt.Fprintf(w, "%s%s:", indent, color.Apply(strings.TrimRight(columns[0], ":"), keyColor))
		} else if len(columns) == 1 {
			fmt.Fprintf(w, "%s%s", indent, color.Apply(columns[0], valColor))
		} else {
			fmt.Fprintf(w, "%s%s", indent, color.Apply(columns[0], keyColor))
		}

		if len(columns) == 1 {
			fmt.Fprint(w, "\n")
			continue
		}

		spacesPos := spacesIndices[0]
		spacesCnt := spacesPos[1] - spacesPos[0]
		fmt.Fprintf(w, "%s%s\n", toSpaces(spacesCnt), color.Apply(columns[1], getColorByValueType(columns[1], dp.DarkBackground)))
	}
}
