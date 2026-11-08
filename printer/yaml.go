package printer

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/internal/bytesutil"
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
	scanner.Buffer(nil, bytesutil.MaxLineLength)
	for scanner.Scan() {
		line := scanner.Text()
		p.printLineAsYAMLFormat(line, w)
	}
	if err := scanner.Err(); err != nil {
		slog.Error("Failed to print YAML output.", "error", err)
	}
}

func (p *YAMLPrinter) printLineAsYAMLFormat(line string, w io.Writer) {
	indentLen := findIndent(line) // can be 0
	indent := line[:indentLen]
	trimmedLine := line[indentLen:]

	if p.nextLineSetsIndent || p.multilineStringIndent > 0 {
		// blank lines inside a block scalar are part of the string,
		// and are commonly written without any indentation
		if trimmedLine == "" {
			fmt.Fprintf(w, "%s\n", line)
			return
		}
	}

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

	if item, ok := strings.CutPrefix(trimmedLine, "- "); ok {
		// - |
		//   multiline string
		if header, num, ok := p.colorizeBlockScalarHeader(item); ok {
			fmt.Fprintf(w, "%s- %s\n", indent, header)
			p.startBlockScalar(indentLen, num)
			return
		}
	}

	if key, val, ok := strings.Cut(trimmedLine, ": "); ok {
		// key: |
		//   multiline string
		if header, num, ok := p.colorizeBlockScalarHeader(val); ok {
			fmt.Fprintf(w, "%s%s: %s\n", indent, p.colorizeYAMLKey(key, indentLen, 2), header)
			// the indentation indicator is relative to the key, which
			// comes after any leading "- " sequence indicators
			p.startBlockScalar(indentLen+countLeadingSequenceIndicators(key)*2, num)
			return
		}

		if quote, afterQuote, ok := stringutil.CutPrefixAny(val, "\"", "'"); ok && !strings.HasSuffix(afterQuote, quote) {
			// key: "value
			// (missing final quote)
			fmt.Fprintf(w, "%s%s: %s%s\n", indent, p.colorizeYAMLKey(key, indentLen, 2), quote, p.Theme.Data.String.Render(afterQuote))
			p.inString = true
			return
		}

		// key: value
		fmt.Fprintf(w, "%s%s: %s\n", indent, p.colorizeYAMLKey(key, indentLen, 2), p.colorizeYAMLValue(val))
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

// colorizeBlockScalarHeader parses a block scalar header, such as "|", ">-",
// "|2", ">+4" or "|2-". The indentation and chomping indicators may come in
// either order. It returns the colorized header and the indentation
// indicator, which is 0 when omitted.
func (p *YAMLPrinter) colorizeBlockScalarHeader(s string) (header string, indent int, ok bool) {
	if s == "" || (s[0] != '|' && s[0] != '>') {
		return "", 0, false
	}
	isChomp := func(s string) bool { return s != "" && (s[0] == '-' || s[0] == '+') }
	isIndent := func(s string) bool { return s != "" && s[0] >= '1' && s[0] <= '9' }

	before, rest := s[:1], s[1:]
	if isChomp(rest) {
		before, rest = s[:2], rest[1:]
	}
	var num, after string
	if isIndent(rest) {
		num, rest = rest[:1], rest[1:]
		if len(before) == 1 && isChomp(rest) {
			after, rest = rest[:1], rest[1:]
		}
	}
	if rest != "" {
		return "", 0, false
	}
	if num == "" {
		return before, 0, true
	}
	return before + p.colorizeYAMLValue(num) + after, int(num[0] - '0'), true
}

func (p *YAMLPrinter) startBlockScalar(parentIndent, indentIndicator int) {
	if indentIndicator == 0 {
		p.nextLineSetsIndent = true
		return
	}
	p.multilineStringIndent = parentIndent + indentIndicator
}

func countLeadingSequenceIndicators(key string) int {
	var n int
	for {
		after, ok := strings.CutPrefix(key, "- ")
		if !ok {
			return n
		}
		key = after
		n++
	}
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
	switch value {
	case "{}", "- {}", "[]", "- []":
		return value
	}

	value, hasLeadingDash := strings.CutPrefix(value, "- ")
	quote, unquotedValue, isDoubleQuoted := stringutil.CutSurroundingAny(value, "\"'")
	color := ColorDataValue(value, p.Theme)

	switch {
	case hasLeadingDash && isDoubleQuoted:
		return fmt.Sprintf(`- %c%s%c`, quote, color.Render(unquotedValue), quote)
	case hasLeadingDash:
		return fmt.Sprintf(`- %s`, color.Render(value))
	case isDoubleQuoted:
		return fmt.Sprintf(`%c%s%c`, quote, color.Render(unquotedValue), quote)
	default:
		return color.Render(value)
	}
}

func (p *YAMLPrinter) colorizeYAMLStringValue(value string) string {
	if before, ok := strings.CutSuffix(value, "\""); ok {
		return fmt.Sprintf(`%s"`, p.Theme.Data.String.Render(before))
	}
	if before, ok := strings.CutSuffix(value, "'"); ok {
		return fmt.Sprintf(`%s'`, p.Theme.Data.String.Render(before))
	}
	return p.Theme.Data.String.Render(value)
}

func (p *YAMLPrinter) isStringClosed(line string) bool {
	return strings.HasSuffix(line, "'") || strings.HasSuffix(line, `"`)
}
