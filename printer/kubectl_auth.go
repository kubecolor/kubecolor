package printer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kubecolor/kubecolor/config"
)

// AuthPrinter is used in "kubectl auth" output
type AuthPrinter struct {
	Theme *config.Theme

	Args []string
	List bool
}

// ensures it implements the interface
var _ Printer = &AuthPrinter{}

func (p *AuthPrinter) Print(r io.Reader, w io.Writer) {
	arg := ""
	if len(p.Args) > 0 {
		arg = p.Args[0]
	}

	switch arg {
	case "can-i":
		p.printCanI(r, w)
	default:
		io.Copy(w, r)
	}
}

func (p *AuthPrinter) printCanI(r io.Reader, w io.Writer) {
	if p.List {
		NewTablePrinter(true, p.Theme, nil).Print(r, w)
	} else {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			switch line {
			case "yes":
				fmt.Fprintf(w, "%s\n", p.Theme.Base.Success.Render(line))
			case "no":
				fmt.Fprintf(w, "%s\n", p.Theme.Base.Warning.Render(line))
			default:
				fmt.Fprintf(w, "%s\n", line)
			}
		}
	}
}
