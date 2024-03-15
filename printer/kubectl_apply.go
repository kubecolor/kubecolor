package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config"
)

type ApplyPrinter struct {
	Theme *config.Theme
}

// kubectl apply
// deployment.apps/foo unchanged
// deployment.apps/bar created
// deployment.apps/quux configured
func (ap *ApplyPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			fmt.Fprintln(w, line)
			continue
		}

		formatted := ap.formatColoredLine(line)
		if formatted == "" {
			fmt.Fprintln(w, ap.Theme.Apply.Fallback.Render(line))
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

func (ap *ApplyPrinter) getColorFor(action string) (config.Color, bool) {
	switch action {
	case "created":
		return ap.Theme.Apply.Created, true
	case "configured":
		return ap.Theme.Apply.Configured, true
	case "unchanged":
		return ap.Theme.Apply.Unchanged, true
	case "(dry run)", "(server dry run)":
		return ap.Theme.Apply.DryRun, true
	default:
		return config.Color{}, false
	}
}
