package printer

import (
	"strconv"
	"strings"

	"github.com/hidetatz/kubecolor/color"
)

// toSpaces returns repeated spaces whose length is n.
func toSpaces(n int) string {
	return strings.Repeat(" ", n)
}

// getColorByKeyIndent returns a color based on the given indent.
// When you want to change key color based on indent depth (e.g. Json, Yaml), use this function
func getColorByKeyIndent(indent int, basicIndentWidth int, dark bool) color.Color {
	switch indent / basicIndentWidth % 2 {
	case 1:
		if dark {
			return color.Cyan
		}
		return color.Black
	default:
		return color.Yellow
	}
}

// getColorByValueType returns a color by value.
// This is intended to be used to colorize any structured data e.g. Json, Yaml.
func getColorByValueType(val string, dark bool) color.Color {
	if val == "null" || val == "<none>" || val == "<unknown>" {
		if dark {
			return NullColorForDark
		}
		return NullColorForLight
	}

	if val == "true" || val == "false" {
		if dark {
			return BoolColorForDark
		}
		return BoolColorForLight
	}

	if _, err := strconv.Atoi(val); err == nil {
		if dark {
			return NumberColorForDark
		}
		return NumberColorForLight
	}

	if dark {
		return StringColorForDark
	}
	return StringColorForLight
}

// getColorsByBackground returns a preset of colors depending on given background color
func getColorsByBackground(dark bool) []color.Color {
	if dark {
		return colorsForDarkBackground
	}

	return colorsForLightBackground
}

// getHeaderColorByBackground returns a defined color for Header (not actual data) by the background color
func getHeaderColorByBackground(dark bool) color.Color {
	if dark {
		return HeaderColorForDark
	}

	return HeaderColorForLight
}

// findIndent returns a length of indent (spaces at left) in the given line
func findIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}
