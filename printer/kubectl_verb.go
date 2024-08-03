package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config"
)

// VerbPrinter is used in change commands like "kubectl apply" output:
//
//	$ kubectl apply
//	deployment.apps/foo unchanged
//	deployment.apps/bar created
//	deployment.apps/quux configured
//
// But it is also used to colorize:
//
// - kubectl create
// - kubectl delete
// - kubectl expose
// - kubectl patch
type VerbPrinter struct {
	DryRunColor   config.Color
	FallbackColor config.Color
	VerbColor     map[string]config.Color
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
		idx := strings.Index(line, verb)
		if idx == -1 || idx == 0 {
			// don't colorize first word
			continue
		}

		if line[idx-1] != ' ' {
			// must be after a space
			continue
		}

		return fmt.Sprintf("%s%s%s", line[:idx], color.Render(verb), line[idx+len(verb):]), true
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
