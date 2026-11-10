package printer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/internal/bytesutil"
)

// RolloutHistoryPrinter is used on "kubectl rollout history" output.
// Each resource starts with its name, followed by a table of revisions:
//
//	deployment.apps/nginx
//	REVISION  CHANGE-CAUSE
//	1         <none>
//	2         <none>
//
// or, with "--revision=N", by the pod template of that revision
// in the same format as "kubectl describe":
//
//	deployment.apps/nginx with revision #2
//	Pod Template:
//	  Labels:       app=nginx
type RolloutHistoryPrinter struct {
	Theme *config.Theme
}

// ensures it implements the interface
var _ Printer = &RolloutHistoryPrinter{}

const rolloutHistoryRevisionSep = " with revision #"

// Print implements [Printer.Print]
func (p *RolloutHistoryPrinter) Print(r io.Reader, w io.Writer) {
	var block bytes.Buffer
	withRevision := false
	flush := func() {
		if block.Len() == 0 {
			return
		}
		if withRevision {
			p.describePrinter().Print(&block, w)
		} else {
			NewTablePrinter(true, p.Theme, nil).Print(&block, w)
		}
		block.Reset()
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(nil, bytesutil.MaxLineLength)
	startOfBlock := true
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			flush()
			fmt.Fprintln(w)
			startOfBlock = true
			continue
		}
		if startOfBlock {
			startOfBlock = false
			if name, revision, ok := parseRolloutHistoryResource(line); ok {
				flush()
				withRevision = revision != ""
				fmt.Fprint(w, p.Theme.Table.Header.Render(name))
				if withRevision {
					fmt.Fprint(w, rolloutHistoryRevisionSep, ColorDataValue(revision, p.Theme).Render(revision))
				} else if trailing := line[len(name):]; trailing != "" {
					fmt.Fprint(w, trailing)
				}
				fmt.Fprintln(w)
				continue
			}
		}
		block.WriteString(line)
		block.WriteByte('\n')
	}
	flush()
	if err := scanner.Err(); err != nil {
		slog.Error("Failed to print rollout history output.", "error", err)
	}
}

func (p *RolloutHistoryPrinter) describePrinter() *DescribePrinter {
	return &DescribePrinter{
		TablePrinter: NewTablePrinter(false, p.Theme, func(_ int, column string) string {
			if colored, ok := ColorStatus(column, p.Theme); ok {
				return colored
			}
			return column
		}),
	}
}

// parseRolloutHistoryResource parses the first line of a resource in
// "kubectl rollout history", such as "deployment.apps/nginx " or
// "deployment.apps/nginx with revision #2".
func parseRolloutHistoryResource(line string) (name, revision string, ok bool) {
	name, revision, _ = strings.Cut(line, rolloutHistoryRevisionSep)
	name = strings.TrimRight(name, " ")
	if name == "" || strings.ContainsAny(name, " \t") || !strings.Contains(name, "/") {
		return "", "", false
	}
	return name, revision, true
}
