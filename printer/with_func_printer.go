package printer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kubecolor/kubecolor/config"
)

// WithFuncPrinter is a printer to print something based on injected logic.
// The function must not be nil, otherwise it panics.
type WithFuncPrinter struct {
	Fn func(line string) config.Color
}

// ensures it implements the interface
var _ Printer = &WithFuncPrinter{}

// Print implements [Printer.Print]
func (p *WithFuncPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		c := p.Fn(line)
		fmt.Fprintf(w, "%s\n", c.Render(line))
	}
}
