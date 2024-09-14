package printer

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/internal/stringutil"
)

// YAMLPrinter is used on "kubectl get -o yaml" output
type YAMLPrinter struct {
	Theme                 *config.Theme
	inString              bool
	multilineStringIndent int
	nextLineSetsIndent    bool
}

// ensures it implements the interface
var _ Printer = &YAMLPrinter{}

// Print implements [Printer.Print]
func (p *YAMLPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		p.printLineAsYAMLFormat(line, w)
	}
}

func (p *YAMLPrinter) printLineAsYAMLFormat(line string, w io.Writer) {
	indentLen := findIndent(line) // can be 0
	indent := line[:indentLen]
	trimmedLine := line[indentLen:]

	if p.nextLineSetsIndent {
		p.multilineStringIndent = indentLen
		p.nextLineSetsIndent = false
	}

	if p.multilineStringIndent > 0 {
		if indentLen >= p.multilineStringIndent {
			fmt.Fprintf(w, "%s%s\n", line[:p.multilineStringIndent], p.Theme.Data.String.Render(line[p.multilineStringIndent:]))
			return
		} else {
			p.multilineStringIndent = 0
		}
	}

	if p.inString {
		// if inString is true, the line must be a part of a string which is broken into several lines
		fmt.Fprintf(w, "%s%s\n", indent, p.colorizeYAMLStringValue(trimmedLine))
		p.inString = !p.isStringClosed(trimmedLine)
		return
	}

	if key, val, ok := strings.Cut(trimmedLine, ": "); ok {
		// key: |
		//   multiline string

		if prefix, after, ok := stringutil.CutPrefixAny(val, "|-", "|+", "|", ">-", ">+", ">"); ok {
			if after == "" {
				fmt.Fprintf(w, "%s%s: %s\n", indent, p.colorizeYAMLKey(key, indentLen, 2), prefix)
				p.nextLineSetsIndent = true
				return
			}
			if num, err := strconv.Atoi(after); err == nil {
				fmt.Fprintf(w, "%s%s: %s%s\n", indent, p.colorizeYAMLKey(key, indentLen, 2), prefix, p.colorizeYAMLValue(after))
				p.multilineStringIndent = indentLen + num
				return
			}
		}

		// key: value
		fmt.Fprintf(w, "%s%s: %s\n", indent, p.colorizeYAMLKey(key, indentLen, 2), p.colorizeYAMLValue(val))
		p.inString = p.isStringOpenedButNotClosed(val)
		return
	}

	// when coming here, the line is just a "key:" or an element of an array
	if strings.HasSuffix(trimmedLine, ":") {
		// key:
		fmt.Fprintf(w, "%s%s\n", indent, p.colorizeYAMLKey(trimmedLine, indentLen, 2))
		return
	}

	fmt.Fprintf(w, "%s%s\n", indent, p.colorizeYAMLValue(trimmedLine))
}

func (p *YAMLPrinter) colorizeYAMLKey(key string, indentCnt, basicWidth int) string {
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

func (p *YAMLPrinter) colorizeYAMLValue(value string) string {
	if value == "{}" {
		return "{}"
	}

	value, hasLeadingDash := strings.CutPrefix(value, "- ")
	unquotedValue, isDoubleQuoted := stringutil.CutSurrounding(value, '"')

	switch {
	case hasLeadingDash && isDoubleQuoted:
		return fmt.Sprintf(`- "%s"`, ColorDataValue(value, p.Theme).Render(unquotedValue))
	case hasLeadingDash:
		return fmt.Sprintf(`- %s`, ColorDataValue(value, p.Theme).Render(value))
	case isDoubleQuoted:
		return fmt.Sprintf(`"%s"`, ColorDataValue(value, p.Theme).Render(unquotedValue))
	default:
		return ColorDataValue(value, p.Theme).Render(value)
	}
}

func (p *YAMLPrinter) colorizeYAMLStringValue(value string) string {
	if withoutQuotes, ok := stringutil.CutSurrounding(value, '"'); ok {
		return fmt.Sprintf(`"%s"`, p.Theme.Data.String.Render(withoutQuotes))
	}
	return p.Theme.Data.String.Render(value)
}

func (p *YAMLPrinter) isStringClosed(line string) bool {
	return strings.HasSuffix(line, "'") || strings.HasSuffix(line, `"`)
}

func (p *YAMLPrinter) isStringOpenedButNotClosed(line string) bool {
	return (strings.HasPrefix(line, "'") && !strings.HasSuffix(line, "'")) ||
		(strings.HasPrefix(line, `"`) && !strings.HasSuffix(line, `"`))
}
