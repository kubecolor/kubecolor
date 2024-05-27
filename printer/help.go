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

// HelpPrinter is a specific printer to print kubectl explain format.
type HelpPrinter struct {
	Theme *config.Theme

	commandBuf              *bytes.Buffer
	lastCommandWasContinued bool
}

var urlRegex = regexp.MustCompile(`\[(https?://[a-zA-Z0-9][-a-zA-Z0-9]*\.[^\]]+)\]|(https?://[a-zA-Z0-9][-a-zA-Z0-9\.@:%_\+~#=/\?]*)`)

func (hp *HelpPrinter) Print(r io.Reader, w io.Writer) {
	scanner := describe.NewScanner(r)
	hp.commandBuf = &bytes.Buffer{}

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
				hp.Theme.Help.Header.Render(string(line.Key)),
				line.Spacing,
				line.Trailing)
			continue
		}

		if (scanner.Path().HasPrefix("Examples") || scanner.Path().HasPrefix("Usage")) &&
			hp.printCommandLine(w, line.String()) {
			continue
		}

		if (scanner.Path().HasPrefix("Options") || scanner.Path().HasPrefix("Flags")) && len(scanner.Path()) == 2 {
			val := string(line.Value)
			fmt.Fprintf(w, "%s%s%s%s%s\n",
				line.Indent,
				hp.Theme.Help.Flag.Render(string(line.Key)),
				line.Spacing,
				getColorByValueType(val, hp.Theme).Render(val),
				line.Trailing)
			continue
		}

		text := string(slices.Concat(line.Key, line.Spacing, line.Value))
		text = hp.colorizeUrls(text)
		if scanner.Path().HasPrefix("Options") {
			fmt.Fprintf(w, "%s%s%s\n",
				line.Indent,
				hp.Theme.Help.FlagDesc.Render(text),
				line.Trailing)
		} else {
			fmt.Fprintf(w, "%s%s%s\n",
				line.Indent,
				hp.Theme.Help.Text.Render(text),
				line.Trailing)
		}
	}
}

func (hp *HelpPrinter) printCommandLine(w io.Writer, line string) bool {
	withoutIndent, ok := strings.CutPrefix(line, "  ")
	if !ok {
		return false
	}
	if withoutIndent == "" {
		return false
	}
	if withoutIndent[0] == '#' {
		fmt.Fprintf(w, "  %s\n", hp.Theme.Shell.Comment.Render(withoutIndent))
		return true
	}

	hp.commandBuf.Reset()

	for pipeIdx, pipe := range strings.Split(withoutIndent, " | ") {
		if pipeIdx > 0 {
			hp.commandBuf.WriteString(" | ")
		}

		// Don't want to use [strings.Fields], as that trims away double-spaces
		fields := strings.Split(pipe, " ")
		for i, field := range fields {
			if i > 0 {
				hp.commandBuf.WriteByte(' ')
			}
			switch {
			case field == "":
				// Empty, do nothing

				// First arg: it's the executable
			case i == 0 && !hp.lastCommandWasContinued,
				// First arg after "kubectl exec --"
				len(fields) > 3 && fields[0] == "kubectl" && fields[1] == "exec" && fields[i-1] == "--":
				hp.commandBuf.WriteString(hp.Theme.Shell.Command.Render(field))
			case field[0] == '-',
				strings.HasPrefix(field, "[(-") && strings.HasSuffix(field, "]"):

				if flag, value, ok := strings.Cut(field, "="); ok {
					hp.commandBuf.WriteString(hp.Theme.Shell.Flag.Render(flag + "="))
					c := getColorByValueType(value, hp.Theme)
					hp.commandBuf.WriteString(c.Render(value))
				} else {
					hp.commandBuf.WriteString(hp.Theme.Shell.Flag.Render(field))
				}
			case isQuoted(field):
				hp.commandBuf.WriteString(hp.Theme.Data.String.Render(field))
			default:
				if flag, value, ok := strings.Cut(field, "="); ok {
					hp.commandBuf.WriteString(hp.Theme.Shell.Arg.Render(flag + "="))
					c := getColorByValueType(value, hp.Theme)
					hp.commandBuf.WriteString(c.Render(value))
				} else {
					hp.commandBuf.WriteString(hp.Theme.Shell.Arg.Render(field))
				}
			}
		}
	}

	fmt.Fprint(w, "  ", hp.commandBuf.String(), "\n")
	hp.lastCommandWasContinued = strings.HasSuffix(withoutIndent, "\\")

	return true
}

func (hp *HelpPrinter) colorizeUrls(s string) string {
	return urlRegex.ReplaceAllStringFunc(s, func(url string) string {
		if url[0] == '[' {
			return fmt.Sprintf("[%s]", hp.Theme.Help.Url.Render(url[1:len(url)-2]))
		}
		return hp.Theme.Help.Url.Render(url)
	})
}

func isQuoted(s string) bool {
	if len(s) < 2 {
		return false
	}
	return (s[0] == '\'' || s[0] == '"') && s[len(s)-1] == s[0]
}
