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

// kubectl apply
// deployment.apps/foo unchanged
// deployment.apps/bar created
// deployment.apps/quux configured
func (ap *ApplyPrinter) Print(r io.Reader, w io.Writer) {
	const (
		applyActionCreated    = "created"
		applyActionConfigured = "configured"
		applyActionUnchanged  = "unchanged"

		dryRunStr = "(dry run)"
	)

	darkColors := map[string]color.Color{
		applyActionCreated:    color.Green,
		applyActionConfigured: color.Yellow,
		applyActionUnchanged:  color.Magenta,
		dryRunStr:             color.Cyan,
	}
	lightColors := map[string]color.Color{
		applyActionCreated:    color.Green,
		applyActionConfigured: color.Yellow,
		applyActionUnchanged:  color.Magenta,
		dryRunStr:             color.Blue,
	}

	colors := func(action string, dark bool) color.Color {
		if dark {
			return darkColors[action]
		}
		return lightColors[action]
	}

	colorize := func(line, action string, dryRun bool, wr io.Writer) {
		if dryRun {
			arg := strings.TrimSuffix(line, fmt.Sprintf(" %s %s", action, dryRunStr))
			fmt.Fprintf(w, "%s %s %s\n",
				arg,
				color.Apply(action, colors(action, ap.DarkBackground)),
				color.Apply(dryRunStr, colors(dryRunStr, ap.DarkBackground)),
			)
			return
		}

		arg := strings.TrimSuffix(line, " "+action)
		fmt.Fprintf(w, "%s %s\n", arg, color.Apply(action, colors(action, ap.DarkBackground)))
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		// on dry run cases, it shows "xxx created (dry run)"
		case strings.HasSuffix(line, fmt.Sprintf(" %s %s", applyActionCreated, dryRunStr)):
			colorize(line, applyActionCreated, true, w)
		case strings.HasSuffix(line, fmt.Sprintf(" %s %s", applyActionConfigured, dryRunStr)):
			colorize(line, applyActionConfigured, true, w)
		case strings.HasSuffix(line, fmt.Sprintf(" %s %s", applyActionUnchanged, dryRunStr)):
			colorize(line, applyActionUnchanged, true, w)

		// not dry run cases, it shows "xxx created"
		case strings.HasSuffix(line, " "+applyActionCreated):
			colorize(line, applyActionCreated, false, w)
		case strings.HasSuffix(line, " "+applyActionConfigured):
			colorize(line, applyActionConfigured, false, w)
		case strings.HasSuffix(line, " "+applyActionUnchanged):
			colorize(line, applyActionUnchanged, false, w)
		default:
			fmt.Fprintf(w, "%s\n", color.Apply(line, color.Green))
		}
	}
}
