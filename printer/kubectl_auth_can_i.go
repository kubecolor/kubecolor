package printer

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/kubecolor/kubecolor/config/color"
	"github.com/kubecolor/kubecolor/internal/bytesutil"
)

// CanIPrinter is used in the literal output of "kubectl auth can-i",
// e.g:
//
//	$ kubectl auth can-i get pods
//	yes
//
//	$ kubectl auth can-i delete pods
//	no
type CanIPrinter struct {
	AllowedColor  color.Color // used on "yes"
	DeniedColor   color.Color // used on "no"
	FallbackColor color.Color // used when the line is neither "yes" nor "no"
}

// ensures it implements the interface
var _ Printer = &CanIPrinter{}

// Print implements [Printer.Print]
func (p *CanIPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(nil, bytesutil.MaxLineLength)
	for scanner.Scan() {
		line := scanner.Text()
		switch strings.TrimSpace(line) {
		case "yes":
			fmt.Fprintln(w, p.AllowedColor.Render(line))
		case "no":
			fmt.Fprintln(w, p.DeniedColor.Render(line))
		default:
			fmt.Fprintln(w, p.FallbackColor.Render(line))
		}
	}
	if err := scanner.Err(); err != nil {
		slog.Error("Failed to print can-i output.", "error", err)
	}
}
