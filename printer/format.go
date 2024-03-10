package printer

import (
	"strings"

	"github.com/kubecolor/kubecolor/config"
)

// toSpaces returns repeated spaces whose length is n.
func toSpaces(n int) string {
	return strings.Repeat(" ", n)
}

// getColorByKeyIndent returns a color based on the given indent.
// When you want to change key color based on indent depth (e.g. Json, Yaml), use this function
func getColorByKeyIndent(indent int, basicIndentWidth int, theme *config.Theme) config.Color {
	switch indent / basicIndentWidth % 2 {
	case 1:
		// TODO: Change to something more appropriate
		return theme.Data.String
	default:
		return theme.Default
	}
}

// getColorByValueType returns a color by value.
// This is intended to be used to colorize any structured data e.g. Json, Yaml.
func getColorByValueType(val string, theme *config.Theme) config.Color {
	switch val {
	case "null", "<none>", "<unknown>", "<unset>", "<nil>":
		return theme.Data.Null
	case "true", "True":
		return theme.Data.True
	case "false", "False":
		return theme.Data.False
	}

	if isOnlyDigits(val) {
		return theme.Data.Number
	}

	return theme.Data.String
}

func isOnlyDigits(s string) bool {
	for _, r := range s {
		if !isDigit(r) {
			return false
		}
	}
	return true
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// findIndent returns a length of indent (spaces at left) in the given line
func findIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}
