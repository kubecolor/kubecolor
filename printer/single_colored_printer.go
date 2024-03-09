package printer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kubecolor/kubecolor/config"
)

// SingleColoredPrinter is a printer to print something in pre-configured color.
type SingleColoredPrinter struct {
	Color config.Color
}

// Print reads r then writes it in w in sp.Color
func (sp *SingleColoredPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Fprintf(w, "%s\n", sp.Color.Render(scanner.Text()))
	}
}
