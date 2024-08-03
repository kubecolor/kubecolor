package printer

import (
	"bytes"
	"io"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/scanner/logscan"
)

// LogsPrinter is used in "kubectl logs" output:
type LogsPrinter struct {
	Theme *config.Theme
}

// ensures it implements the interface
var _ Printer = &LogsPrinter{}

// Print implements [Printer.Print]
func (p *LogsPrinter) Print(r io.Reader, w io.Writer) {
	scanner := logscan.NewScanner(r)

	// Buffer the lines so we can write them to the io.Writer all at once
	var lineBuffer bytes.Buffer
	var keyIndex int

	for scanner.Scan() {
		token := scanner.Token()

		switch token.Kind {
		case logscan.KindKey:
			var color config.Color
			if len(p.Theme.Logs.Key) > 0 {
				color = p.Theme.Logs.Key[keyIndex%len(p.Theme.Logs.Key)]
			}
			lineBuffer.WriteString(color.Render(token.Text))
			keyIndex++
		case logscan.KindValue:
			lineBuffer.WriteString(ColorDataValue(token.Text, p.Theme).Render(token.Text))
		case logscan.KindQuote:
			lineBuffer.WriteString(p.Theme.Data.String.Render(token.Text))

		case logscan.KindSeverityTrace:
			lineBuffer.WriteString(p.Theme.Logs.Severity.Trace.Render(token.Text))
		case logscan.KindSeverityDebug:
			lineBuffer.WriteString(p.Theme.Logs.Severity.Debug.Render(token.Text))
		case logscan.KindSeverityInfo:
			lineBuffer.WriteString(p.Theme.Logs.Severity.Info.Render(token.Text))
		case logscan.KindSeverityWarn:
			lineBuffer.WriteString(p.Theme.Logs.Severity.Warn.Render(token.Text))
		case logscan.KindSeverityError:
			lineBuffer.WriteString(p.Theme.Logs.Severity.Error.Render(token.Text))
		case logscan.KindSeverityFatal:
			lineBuffer.WriteString(p.Theme.Logs.Severity.Fatal.Render(token.Text))
		case logscan.KindSeverityPanic:
			lineBuffer.WriteString(p.Theme.Logs.Severity.Panic.Render(token.Text))

		case logscan.KindNewline:
			lineBuffer.WriteByte('\n')
			lineBuffer.WriteTo(w)
			lineBuffer.Reset()
			keyIndex = 0

		default:
			lineBuffer.WriteString(token.Text)
		}
	}
}
