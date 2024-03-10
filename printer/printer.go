package printer

import (
	"io"
)

// Printer can print something.
// It reads data from r, then write them in w.
type Printer interface {
	Print(r io.Reader, w io.Writer)
}
