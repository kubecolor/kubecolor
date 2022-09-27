package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/color"
)

// ExplainPrinter is a specific printer to print kubectl explain format.
type ExplainPrinter struct {
	DarkBackground bool
	Recursive      bool

	renderingFields bool
}

func (ep *ExplainPrinter) Print(r io.Reader, w io.Writer) {
	// https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/kubectl/pkg/explain/model_printer.go#L24-L30
	descriptionIndentLevel := 5

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			fmt.Fprintln(w)
			continue
		}

		if ep.renderingFields {
			ep.printField(line, w)
			continue
		}

		indentCnt := findIndent(line)
		switch indentCnt {
		case 0:
			ep.printKeyVal(line, w)
		case descriptionIndentLevel:
			ep.printDescription(line, w)
		}

		if line == "FIELDS:" {
			ep.renderingFields = true
		}
	}
}

func (ep *ExplainPrinter) printKeyVal(line string, w io.Writer) {
	var key, val string
	keyAndVal := spaces.Split(line, 2)
	if len(keyAndVal) == 1 {
		key = line
	} else {
		key, val = keyAndVal[0], keyAndVal[1]
	}

	key = strings.TrimRight(key, ":")

	key = color.Apply(key, getColorByKeyIndent(0, 2, ep.DarkBackground))
	if val != "" {
		val = color.Apply(val, getColorByValueType(val, ep.DarkBackground))
	}

	spacesIndices := spaces.FindAllStringIndex(line, -1)
	spacesBetweenKeyAndVal := ""
	if len(spacesIndices) > 0 {
		spacesBetweenKeyAndVal = toSpaces(spacesIndices[0][1] - spacesIndices[0][0])
	}

	fmt.Fprintf(w, "%s:%s%s\n", key, spacesBetweenKeyAndVal, val)
}

func (ep *ExplainPrinter) printDescription(line string, w io.Writer) {
	fmt.Fprintf(w, "%s%s\n", toSpaces(5), color.Apply(strings.TrimLeft(line, " "), getColorByValueType(line, ep.DarkBackground)))

}

func (ep *ExplainPrinter) printField(line string, w io.Writer) {
	if ep.Recursive {
		ep.printKeyAndType(line, w)
		return
	}

	indentCnt := findIndent(line)
	if indentCnt == 3 {
		ep.printKeyAndType(line, w)
		return
	}

	ep.printDescription(line, w)
}

func (ep *ExplainPrinter) printKeyAndType(line string, w io.Writer) {
	indentCnt := findIndent(line)
	line = strings.TrimLeft(line, " ")

	keyAndVal := singleOrMultipleSpaces.Split(line, 2) // spaces between key and type can be only 1 space
	key, val := keyAndVal[0], keyAndVal[1]

	val = strings.TrimLeft(strings.TrimRight(val, ">"), "<")
	key = color.Apply(key, getColorByKeyIndent(indentCnt, 2, ep.DarkBackground))
	val = color.Apply(val, getColorByValueType(line, ep.DarkBackground))

	// I don't know why but kubectl explain uses \t as delimiter
	fmt.Fprintf(w, "%s%s\t<%s>\n", toSpaces(indentCnt), key, val)
}
