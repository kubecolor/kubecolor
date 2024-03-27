package printer

import (
	"regexp"
	"strings"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/internal/stringutil"
)

// toSpaces returns repeated spaces whose length is n.
func toSpaces(n int) string {
	return strings.Repeat(" ", n)
}

// getColorByKeyIndent returns a color based on the given indent.
// When you want to change key color based on indent depth (e.g. Json, Yaml), use this function
func getColorByKeyIndent(indent int, basicIndentWidth int, colors config.ColorSlice) config.Color {
	if len(colors) == 0 {
		return config.Color{}
	}
	return colors[indent/basicIndentWidth%len(colors)]
}

var (
	isQuantityRegex = regexp.MustCompile(`^[\+\-]?(?:\d+|\.\d+|\d+\.|\d+\.\d+)?(?:m|[kMGTPE]i?)$`)
)

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

	// Ints: 123
	if stringutil.IsOnlyDigits(val) {
		return theme.Data.Number
	}

	// Floats: 123.456
	if left, right, ok := strings.Cut(val, "."); ok {
		if stringutil.IsOnlyDigits(left) && stringutil.IsOnlyDigits(right) {
			return theme.Data.Number
		}
	}

	// Quantity: 100m, 5Gi
	if isQuantityRegex.MatchString(val) {
		return theme.Data.Quantity
	}

	// Duration: 15m10s
	if _, ok := stringutil.ParseHumanDuration(val); ok {
		return theme.Data.Duration
	}

	return theme.Data.String
}

// findIndent returns a length of indent (spaces at left) in the given line
func findIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}
