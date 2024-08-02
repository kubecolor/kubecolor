package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config"
)

type JSONPrinter struct {
	Theme *config.Theme
}

func (p *JSONPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		printLineAsJsonFormat(line, w, p.Theme)
	}
}

func printLineAsJsonFormat(line string, w io.Writer, theme *config.Theme) {
	indentCnt := findIndent(line)
	indent := toSpaces(indentCnt)
	trimmedLine := strings.TrimLeft(line, " ")

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
	split := strings.SplitN(trimmedLine, ": ", 2) // if key contains ": " this works in a wrong way but it's unlikely to happen

	if len(split) == 1 {
		// when coming here, it will be a value in an array
		fmt.Fprintf(w, "%s%s\n", indent, toColorizedJsonValue(split[0], theme))
		return
	}

	key := split[0]
	val := split[1]

	fmt.Fprintf(w, "%s%s: %s\n", indent, toColorizedJsonKey(key, indentCnt, 4, theme), toColorizedJsonValue(val, theme))
}

// toColorizedJsonKey returns colored json key
func toColorizedJsonKey(key string, indentCnt, basicWidth int, theme *config.Theme) string {
	hasColon := strings.HasSuffix(key, ":")
	// remove colon and double quotations although they might not exist actually
	key = strings.TrimRight(key, ":")
	doubleQuoteTrimmed := strings.TrimRight(strings.TrimLeft(key, `"`), `"`)

	format := `"%s"`
	if hasColon {
		format += ":"
	}

	return fmt.Sprintf(format, ColorDataKey(indentCnt, basicWidth, theme.Data.Key).Render(doubleQuoteTrimmed))
}

// toColorizedJsonValue returns colored json value.
// This function checks it trailing comma and double quotation exist
// then colorize the given value considering them.
func toColorizedJsonValue(value string, theme *config.Theme) string {
	if value == "{" {
		return "{"
	}

	if value == "[" {
		return "["
	}

	if value == "{}," {
		return "{},"
	}

	if value == "{}" {
		return "{}"
	}

	hasComma := strings.HasSuffix(value, ",")
	// remove comma and double quotations although they might not exist actually
	value = strings.TrimRight(value, ",")

	isString := strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)
	doubleQuoteTrimmedValue := strings.TrimRight(strings.TrimLeft(value, `"`), `"`)

	var format string
	switch {
	case hasComma && isString:
		format = `"%s",`
	case hasComma:
		format = `%s,`
	case isString:
		format = `"%s"`
	default:
		format = `%s`
	}

	return fmt.Sprintf(format, ColorDataValue(value, theme).Render(doubleQuoteTrimmedValue))
}
