package printer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kubecolor/kubecolor/config"
)

// WithFuncPrinter is a printer to print something based on injected logic.
type WithFuncPrinter struct {
	Fn func(line string) config.Color
}

// Print reads r then writes it in w but its color is decided by
// pre-injected function.
// The function must not be nil, otherwise it panics.
func (wp *WithFuncPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		c := wp.Fn(line)
		fmt.Fprintf(w, "%s\n", c.Render(line))
	}
}
