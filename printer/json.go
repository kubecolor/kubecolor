package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/internal/stringutil"
)

type JSONPrinter struct {
	Theme *config.Theme
}

func (p *JSONPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		p.printLineAsJsonFormat(line, w)
	}
}

func (p *JSONPrinter) printLineAsJsonFormat(line string, w io.Writer) {
	indentLen := findIndent(line)
	indent := line[:indentLen]
	trimmedLine := line[indentLen:]

	if strings.HasPrefix(trimmedLine, "{") ||
		strings.HasPrefix(trimmedLine, "}") ||
		strings.HasPrefix(trimmedLine, "]") {
		// when coming here, it must not be starting with key.
		// that patterns are:
		// {
		// }
		// },
		// ]
		// ],
		// note: it must not be "[" because it will be always after key
		// in this case, just write it without color
		// fmt.Fprintf(w, "%s", toSpaces(indentCnt))
		fmt.Fprintf(w, "%s%s", indent, trimmedLine)
		fmt.Fprintf(w, "\n")
		return
	}

	// when coming here, the line must be one of:
	// "key": {
	// "key": [
	// "key": value
	// "key": value,
	// value,
	// value

	if key, val, ok := strings.Cut(trimmedLine, ": "); ok {
		fmt.Fprintf(w, "%s%s: %s\n", indent, p.colorizeJSONKey(key, indentLen, 4), p.colorizeJSONValue(val))
		return
	}

	// when coming here, it will be a value in an array
	fmt.Fprintf(w, "%s%s\n", indent, p.colorizeJSONValue(trimmedLine))
}

// colorizeJSONKey returns colored json key
func (p *JSONPrinter) colorizeJSONKey(key string, indentCnt, basicWidth int) string {
	key, isDoubleQuoted := stringutil.CutSurrounding(key, '"')
	color := ColorDataKey(indentCnt, basicWidth, p.Theme.Data.Key)

	switch {
	case isDoubleQuoted:
		return fmt.Sprintf(`"%s"`, color.Render(key))
	default:
		return color.Render(key)
	}
}

// colorizeJSONValue returns colored json value.
// This function checks it trailing comma and double quotation exist
// then colorize the given value considering them.
func (p *JSONPrinter) colorizeJSONValue(value string) string {
	switch value {
	case "{", "[", "{},", "{}":
		return value
	}

	value, hasComma := strings.CutSuffix(value, ",")
	unquotedValue, isDoubleQuoted := stringutil.CutSurrounding(value, '"')
	color := ColorDataValue(value, p.Theme)

	switch {
	case hasComma && isDoubleQuoted:
		return fmt.Sprintf(`"%s",`, color.Render(unquotedValue))
	case hasComma:
		return fmt.Sprintf(`%s,`, color.Render(value))
	case isDoubleQuoted:
		return fmt.Sprintf(`"%s"`, color.Render(unquotedValue))
	default:
		return color.Render(value)
	}
}
