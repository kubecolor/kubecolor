package printer

import (
	"strings"

	"github.com/kubecolor/kubecolor/color"
)

// toSpaces returns repeated spaces whose length is n.
func toSpaces(n int) string {
	return strings.Repeat(" ", n)
}

// getColorByKeyIndent returns a color based on the given indent.
// When you want to change key color based on indent depth (e.g. Json, Yaml), use this function
func getColorByKeyIndent(indent int, basicIndentWidth int, schema ColorSchema) color.Color {
	switch indent / basicIndentWidth % 2 {
	case 1:
		return schema.StringColor
	default:
		return schema.DefaultColor
	}
}

// getColorByValueType returns a color by value.
// This is intended to be used to colorize any structured data e.g. Json, Yaml.
func getColorByValueType(val string, schema ColorSchema) color.Color {
	switch val {
	case "null", "<none>", "<unknown>", "<unset>":
		return schema.NullColor
	case "true", "True":
		return schema.TrueColor
	case "false", "False":
		return schema.FalseColor
	}

	if isOnlyDigits(val) {
		return schema.NumberColor
	}

	return schema.StringColor
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

// getColorsByBackground returns a preset of colors depending on given background color
// func getColorsByBackground(dark bool) []color.Color {
// 	if dark {
// 		return colorsForDarkBackground
// 	}

// 	return colorsForLightBackground
// }

// getHeaderColorByBackground returns a defined color for Header (not actual data) by the background color
// func getHeaderColorByBackground(dark bool) color.Color {
// 	if dark {
// 		return HeaderColorForDark
// 	}

// 	return HeaderColorForLight
// }

// findIndent returns a length of indent (spaces at left) in the given line
func findIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}
