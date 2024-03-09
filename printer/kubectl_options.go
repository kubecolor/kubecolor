package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config"
)

type OptionsPrinter struct {
	Theme *config.Theme
}

func (op *OptionsPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	isFirstLine := true
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			fmt.Fprintln(w)
			continue
		}

		if isFirstLine {
			fmt.Fprintf(w, "%s\n", op.Theme.String.Render(line))
			isFirstLine = false
			continue
		}

		indentCnt := findIndent(line)
		indent := toSpaces(indentCnt)
		trimmedLine := strings.TrimLeft(line, " ")

		splitted := strings.SplitN(trimmedLine, ": ", 2)
		key, val := splitted[0], splitted[1]

		fmt.Fprintf(w, "%s%s: %s\n", indent, getColorByKeyIndent(0, 2, op.Theme).Render(key), getColorByValueType(val, op.Theme).Render(val))
	}
}
