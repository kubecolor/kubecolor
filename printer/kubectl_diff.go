package printer

import (
  "bufio"
  "fmt"
  "github.com/kubecolor/kubecolor/config"
  "io"
  "regexp"
)

type DiffPrinter struct {
  Theme *config.Theme
}

func (p *DiffPrinter) Print(r io.Reader, w io.Writer) {
  scanner := bufio.NewScanner(r)
  for scanner.Scan() {
    line := scanner.Text()
    if line == "" {
      fmt.Fprintln(w, line)
      continue
    }
    parsedLine := p.parseLine(line)

    fmt.Fprintln(w, parsedLine)
  }
}

func (p *DiffPrinter) parseLine(line string) string {
  theme := p.Theme.Diff
  switch {
  // All only lines, that starts with a single "+"and anything afterward
  // e.g. + foo
  // not allowed: ++ foo
  case regexp.MustCompile(`^\+[^+]+.*$`).MatchString(line):
    return theme.Added.Render(line)
  // All only lines, that starts with a single "-" and anything afterward
  // e.g. - foo
  // not allowed: -- foo
  case regexp.MustCompile(`^-[^-]+.*$`).MatchString(line):
    return theme.Removed.Render(line)
  default:
    return theme.Unchanged.Render(line)
  }
}
