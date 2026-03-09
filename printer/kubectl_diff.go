package printer

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/kubecolor/kubecolor/config"
)

type DiffPrinter struct {
	Theme *config.Theme
}

func (p *DiffPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(nil, bufio.MaxScanTokenSize)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			fmt.Fprintln(w, line)
			continue
		}
		parsedLine := p.parseLine(line)

		fmt.Fprintln(w, parsedLine)
	}
	if err := scanner.Err(); err != nil {
		slog.Error("Failed to print diff output.", "error", err)
	}
}

func (p *DiffPrinter) parseLine(line string) string {
	theme := p.Theme.Diff
	switch {
	case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
		return theme.Added.Render(line)
	case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
		return theme.Removed.Render(line)
	default:
		return theme.Unchanged.Render(line)
	}
}
