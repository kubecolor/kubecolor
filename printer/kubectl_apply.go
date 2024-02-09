package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/color"
)

type ApplyPrinter struct {
	Theme *color.Theme
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
			fmt.Fprintln(w, color.Apply(line, color.Green))
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
		return fmt.Sprintf("%s %s", resource, color.Apply(action, actionColor))
	}
	dryRun := strings.TrimPrefix(after, " ")
	dryRunColor, dryRunOk := ap.getColorFor(dryRun)
	if !dryRunOk {
		return ""
	}
	return fmt.Sprintf("%s %s %s", resource, color.Apply(action, actionColor), color.Apply(dryRun, dryRunColor))
}

func (ap *ApplyPrinter) getColorFor(action string) (color.Color, bool) {
	switch action {
	case "created":
		return ap.Theme.Apply.CreatedColor, true
	case "configured":
		return ap.Theme.Apply.ConfiguredColor, true
	case "unchanged":
		return ap.Theme.Apply.UnchangedColor, true
	case "(dry run)", "(server dry run)":
		return ap.Theme.Apply.DryRunColor, true
	default:
		return 0, false
	}
}
