package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config"
)

// YAMLPrinter is used on "kubectl get -o yaml" output
type YAMLPrinter struct {
	Theme    *config.Theme
	inString bool
}

// ensures it implements the interface
var _ Printer = &YAMLPrinter{}

// Print implements [Printer.Print]
func (p *YAMLPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		p.printLineAsYamlFormat(line, w)
	}
}

func (p *YAMLPrinter) printLineAsYamlFormat(line string, w io.Writer) {
	indentCnt := findIndent(line) // can be 0
	indent := toSpaces(indentCnt) // so, can be empty
	trimmedLine := strings.TrimLeft(line, " ")

	if p.inString {
		// if inString is true, the line must be a part of a string which is broken into several lines
		fmt.Fprintf(w, "%s%s\n", indent, p.toColorizedStringValue(trimmedLine))
		p.inString = !p.isStringClosed(trimmedLine)
		return
	}

	split := strings.SplitN(trimmedLine, ": ", 2) // assuming key does not contain ": " while value might do

	if len(split) == 2 {
		// key: value
		key, val := split[0], split[1]
		fmt.Fprintf(w, "%s%s: %s\n", indent, p.toColorizedYamlKey(key, indentCnt, 2), p.toColorizedYamlValue(val))
		p.inString = p.isStringOpenedButNotClosed(val)
		return
	}

	// when coming here, the line is just a "key:" or an element of an array
	if strings.HasSuffix(split[0], ":") {
		// key:
		fmt.Fprintf(w, "%s%s\n", indent, p.toColorizedYamlKey(split[0], indentCnt, 2))
		return
	}

	fmt.Fprintf(w, "%s%s\n", indent, p.toColorizedYamlValue(split[0]))
}

func (p *YAMLPrinter) toColorizedYamlKey(key string, indentCnt, basicWidth int) string {
	hasColon := strings.HasSuffix(key, ":")
	hasLeadingDash := strings.HasPrefix(key, "- ")
	key = strings.TrimSuffix(key, ":")
	key = strings.TrimPrefix(key, "- ")

	format := "%s"
	if hasColon {
		format += ":"
	}

	if hasLeadingDash {
		format = "- " + format
		indentCnt += 2
	}

	return fmt.Sprintf(format, ColorDataKey(indentCnt, basicWidth, p.Theme.Data.Key).Render(key))
}

func (p *YAMLPrinter) toColorizedYamlValue(value string) string {
	if value == "{}" {
		return "{}"
	}

	hasLeadingDash := strings.HasPrefix(value, "- ")
	value = strings.TrimPrefix(value, "- ")

	isDoubleQuoted := strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)
	trimmedValue := strings.TrimSuffix(strings.TrimPrefix(value, `"`), `"`)

	var format string
	switch {
	case hasLeadingDash && isDoubleQuoted:
		format = `- "%s"`
	case hasLeadingDash:
		format = "- %s"
	case isDoubleQuoted:
		format = `"%s"`
	default:
		format = "%s"
	}

	return fmt.Sprintf(format, ColorDataValue(value, p.Theme).Render(trimmedValue))
}

func (p *YAMLPrinter) toColorizedStringValue(value string) string {

	isDoubleQuoted := strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)
	trimmedValue := strings.TrimRight(strings.TrimLeft(value, `"`), `"`)

	var format string
	switch {
	case isDoubleQuoted:
		format = `"%s"`
	default:
		format = "%s"
	}
	return fmt.Sprintf(format, p.Theme.Data.String.Render(trimmedValue))
}

func (p *YAMLPrinter) isStringClosed(line string) bool {
	return strings.HasSuffix(line, "'") || strings.HasSuffix(line, `"`)
}

func (p *YAMLPrinter) isStringOpenedButNotClosed(line string) bool {
	return (strings.HasPrefix(line, "'") && !strings.HasSuffix(line, "'")) ||
		(strings.HasPrefix(line, `"`) && !strings.HasSuffix(line, `"`))
}
