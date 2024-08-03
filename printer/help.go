package printer

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/scanner/describe"
)

// HelpPrinter is used on "kubectl --help" output
type HelpPrinter struct {
	Theme *config.Theme

	commandBuf              *bytes.Buffer
	lastCommandWasContinued bool
}

// ensures it implements the interface
var _ Printer = &HelpPrinter{}

var urlRegex = regexp.MustCompile(`\[(https?://[a-zA-Z0-9][-a-zA-Z0-9]*\.[^\]]+)\]|(https?://[a-zA-Z0-9][-a-zA-Z0-9\.@:%_\+~#=/\?]*)`)

// Print implements [Printer.Print]
func (p *HelpPrinter) Print(r io.Reader, w io.Writer) {
	scanner := describe.NewScanner(r)
	p.commandBuf = &bytes.Buffer{}

	for scanner.Scan() {
		line := scanner.Line()

		if line.IsZero() {
			fmt.Fprintln(w)
			continue
		}

		if len(line.Value) == 0 &&
			len(scanner.Path()) == 1 &&
			strings.HasSuffix(line.String(), ":") {
			fmt.Fprintf(w, "%s%s%s%s\n",
				line.Indent,
				p.Theme.Help.Header.Render(string(line.Key)),
				line.Spacing,
				line.Trailing)
			continue
		}

		if (scanner.Path().HasPrefix("Examples") || scanner.Path().HasPrefix("Usage")) &&
			p.printCommandLine(w, line.String()) {
			continue
		}

		if (scanner.Path().HasPrefix("Options") || scanner.Path().HasPrefix("Flags")) && len(scanner.Path()) == 2 {
			val := string(line.Value)
			fmt.Fprintf(w, "%s%s%s%s%s\n",
				line.Indent,
				p.Theme.Help.Flag.Render(string(line.Key)),
				line.Spacing,
				ColorDataValue(val, p.Theme).Render(val),
				line.Trailing)
			continue
		}

		text := string(slices.Concat(line.Key, line.Spacing, line.Value))
		text = p.colorizeUrls(text)
		if scanner.Path().HasPrefix("Options") {
			fmt.Fprintf(w, "%s%s%s\n",
				line.Indent,
				p.Theme.Help.FlagDesc.Render(text),
				line.Trailing)
		} else {
			fmt.Fprintf(w, "%s%s%s\n",
				line.Indent,
				p.Theme.Help.Text.Render(text),
				line.Trailing)
		}
	}
}

func (p *HelpPrinter) printCommandLine(w io.Writer, line string) bool {
	withoutIndent, ok := strings.CutPrefix(line, "  ")
	if !ok {
		return false
	}
	if withoutIndent == "" {
		return false
	}
	if withoutIndent[0] == '#' {
		fmt.Fprintf(w, "  %s\n", p.Theme.Shell.Comment.Render(withoutIndent))
		return true
	}

	p.commandBuf.Reset()

	for pipeIdx, pipe := range strings.Split(withoutIndent, " | ") {
		if pipeIdx > 0 {
			p.commandBuf.WriteString(" | ")
		}

		// Don't want to use [strings.Fields], as that trims away double-spaces
		fields := strings.Split(pipe, " ")
		for i, field := range fields {
			if i > 0 {
				p.commandBuf.WriteByte(' ')
			}
			switch {
			case field == "":
				// Empty, do nothing

				// First arg: it's the executable
			case i == 0 && !p.lastCommandWasContinued,
				// First arg after "kubectl exec --"
				len(fields) > 3 && fields[0] == "kubectl" && fields[1] == "exec" && fields[i-1] == "--":
				p.commandBuf.WriteString(p.Theme.Shell.Command.Render(field))
			case field[0] == '-',
				strings.HasPrefix(field, "[(-") && strings.HasSuffix(field, "]"):

				if flag, value, ok := strings.Cut(field, "="); ok {
					p.commandBuf.WriteString(p.Theme.Shell.Flag.Render(flag + "="))
					c := ColorDataValue(value, p.Theme)
					p.commandBuf.WriteString(c.Render(value))
				} else {
					p.commandBuf.WriteString(p.Theme.Shell.Flag.Render(field))
				}
			case isQuoted(field):
				p.commandBuf.WriteString(p.Theme.Data.String.Render(field))
			default:
				if flag, value, ok := strings.Cut(field, "="); ok {
					p.commandBuf.WriteString(p.Theme.Shell.Arg.Render(flag + "="))
					c := ColorDataValue(value, p.Theme)
					p.commandBuf.WriteString(c.Render(value))
				} else {
					p.commandBuf.WriteString(p.Theme.Shell.Arg.Render(field))
				}
			}
		}
	}

	fmt.Fprint(w, "  ", p.commandBuf.String(), "\n")
	p.lastCommandWasContinued = strings.HasSuffix(withoutIndent, "\\")

	return true
}

func (p *HelpPrinter) colorizeUrls(s string) string {
	return urlRegex.ReplaceAllStringFunc(s, func(url string) string {
		if url[0] == '[' {
			return fmt.Sprintf("[%s]", p.Theme.Help.Url.Render(url[1:len(url)-2]))
		}
		return p.Theme.Help.Url.Render(url)
	})
}

func isQuoted(s string) bool {
	if len(s) < 2 {
		return false
	}
	return (s[0] == '\'' || s[0] == '"') && s[len(s)-1] == s[0]
}
