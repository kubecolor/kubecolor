package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config"
)

// ApplyPrinter is used in "kubectl apply" output:
//
//	kubectl apply
//	deployment.apps/foo unchanged
//	deployment.apps/bar created
//	deployment.apps/quux configured
type ApplyPrinter struct {
	Theme *config.Theme
}

// ensures it implements the interface
var _ Printer = &ApplyPrinter{}

// Print implements [Printer.Print]
func (p *ApplyPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			fmt.Fprintln(w, line)
			continue
		}

		formatted := p.formatColoredLine(line)
		if formatted == "" {
			fmt.Fprintln(w, p.Theme.Apply.Fallback.Render(line))
			continue
		}

		fmt.Fprintln(w, formatted)
	}
}

func (ap *ApplyPrinter) formatColoredLine(line string) string {
	resource, after, hasAction := strings.Cut(line, " ")
	if !hasAction {
		return ""
	}
	action, after, hasDryRun := strings.Cut(after, " ")
	actionColor, actionOk := ap.getColorFor(action)
	if !actionOk {
		return ""
	}
	if !hasDryRun {
		return fmt.Sprintf("%s %s", resource, actionColor.Render(action))
	}
	dryRun := strings.TrimPrefix(after, " ")
	dryRunColor, dryRunOk := ap.getColorFor(dryRun)
	if !dryRunOk {
		return ""
	}
	return fmt.Sprintf("%s %s %s", resource, actionColor.Render(action), dryRunColor.Render(dryRun))
}

func (p *ApplyPrinter) getColorFor(action string) (config.Color, bool) {
	switch action {
	case "created":
		return p.Theme.Apply.Created, true
	case "configured":
		return p.Theme.Apply.Configured, true
	case "unchanged":
		return p.Theme.Apply.Unchanged, true
	case "(dry run)", "(server dry run)":
		return p.Theme.Apply.DryRun, true
	default:
		return config.Color{}, false
	}
}
