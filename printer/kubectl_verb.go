package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config/color"
)

// VerbPrinter is used in change commands like "kubectl apply" output:
//
//	$ kubectl apply
//	deployment.apps/foo unchanged
//	deployment.apps/bar created
//	deployment.apps/quux configured
//
// But it is also used to colorize commands like:
//
// - kubectl create
// - kubectl delete
// - kubectl drain
// - kubectl expose
// - kubectl patch
// - kubectl uncordon
type VerbPrinter struct {
	// VerbColor is used for verbs at the end of the line (followed by optional "dry-run"), e.g:
	// 	 pod/nginx-28729634-nh2vc evicted
	// 	 pod "nginx-28729634-nh2vc" deleted (dry run)
	VerbColor map[string]color.Color

	// PrefixVerbColor is used for verbs that are prefixes instead of suffixes, e.g:
	// 	 evicting pod nginx/nginx-28729634-nh2vc
	PrefixVerbColor map[string]color.Color

	// DryRunColor is used on the "(dry run)" or "(server dry run)" suffix of a line.
	DryRunColor color.Color

	// FallbackColor is used when no verbs has matched on the output line.
	FallbackColor color.Color
}

// ensures it implements the interface
var _ Printer = &VerbPrinter{}

const (
	dryRunClientSuffix = "(dry run)"
	dryRunServerSuffix = "(server dry run)"
)

// Print implements [Printer.Print]
func (p *VerbPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			fmt.Fprintln(w, line)
			continue
		}

		colored, isColored := p.colorizeVerb(line)
		if isColored {
			colored = p.colorizeDryRun(colored)
		} else {
			colored = p.FallbackColor.Render(line)
		}

		fmt.Fprintln(w, colored)
	}
}

func (p *VerbPrinter) colorizeVerb(line string) (string, bool) {
	for verb, color := range p.VerbColor {
		idx := strings.LastIndex(line, verb)
		if idx == -1 || idx == 0 {
			// don't colorize first word
			continue
		}
		before := line[:idx]
		after := line[idx+len(verb):]

		if !strings.HasSuffix(before, " ") {
			// must be after a space
			continue
		}

		return fmt.Sprintf("%s%s%s", before, color.Render(verb), after), true
	}

	for verbPrefix, color := range p.PrefixVerbColor {
		after, ok := strings.CutPrefix(line, verbPrefix)
		if !ok {
			continue
		}

		if !strings.HasPrefix(after, " ") {
			// must be before a space
			continue
		}

		return color.Render(verbPrefix) + after, true
	}

	return line, false
}

func (p *VerbPrinter) colorizeDryRun(line string) string {
	if before, ok := strings.CutSuffix(line, dryRunClientSuffix); ok {
		return before + p.DryRunColor.Render(dryRunClientSuffix)
	}
	if before, ok := strings.CutSuffix(line, dryRunServerSuffix); ok {
		return before + p.DryRunColor.Render(dryRunServerSuffix)
	}
	return line
}
