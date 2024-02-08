package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/color"
)

type ApplyPrinter struct {
	DarkBackground bool
}

const (
	applyActionCreatedStr    = "created"
	applyActionConfiguredStr = "configured"
	applyActionUnchangedStr  = "unchanged"

	applyDryRunStr       = "(dry run)"
	applyDryRunServerStr = "(server dry run)"
)

var applyDarkColors = map[string]color.Color{
	applyActionCreatedStr:    color.Green,
	applyActionConfiguredStr: color.Yellow,
	applyActionUnchangedStr:  color.Magenta,
	applyDryRunStr:           color.Cyan,
	applyDryRunServerStr:     color.Cyan,
}
var applyLightColors = map[string]color.Color{
	applyActionCreatedStr:    color.Green,
	applyActionConfiguredStr: color.Yellow,
	applyActionUnchangedStr:  color.Magenta,
	applyDryRunStr:           color.Blue,
	applyDryRunServerStr:     color.Cyan,
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
	if ap.DarkBackground {
		c, ok := applyDarkColors[action]
		return c, ok
	}
	c, ok := applyLightColors[action]
	return c, ok
}
