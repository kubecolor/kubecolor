package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config"
)

// StderrPrinter is a used on stderr input.
type StderrPrinter struct {
	Theme *config.Theme
}

// ensures it implements the interface
var _ Printer = &StderrPrinter{}

// Print implements [Printer.Print]
func (p *StderrPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)

	logsPrinter := LogsPrinter{
		Theme: p.Theme,
	}
	logsPrinterReader := strings.NewReader("")

	for scanner.Scan() {
		line := scanner.Text()
		if s, ok := p.formatLine(line); ok {
			fmt.Fprintln(w, s)
			continue
		}

		logsPrinterReader.Reset(line)
		logsPrinter.Print(logsPrinterReader, w)
	}
}

func (p *StderrPrinter) formatLine(line string) (string, bool) {
	if strings.HasPrefix(strings.ToLower(line), "error") {
		return p.Theme.Stderr.Error.Render(line), true
	}

	if after, ok := strings.CutPrefix(line, "No resources found"); ok {
		if after == "" {
			return p.Theme.Stderr.NoneFound.Render(line), true
		}
		if afterIn, ok := strings.CutPrefix(after, " in "); ok {
			if ns, ok := strings.CutSuffix(afterIn, " namespace."); ok {
				return p.Theme.Stderr.NoneFound.Render(fmt.Sprintf(
					"No resources found in %s namespace.",
					p.Theme.Stderr.NoneFoundNamespace.Render(ns),
				)), true
			}
		}
	}

	return "", false
}
