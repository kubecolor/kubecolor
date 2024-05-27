package config

import (
	"fmt"

	"github.com/gookit/color"
)

func ConvertColorCode(c ColorCode) string {
	switch c := c.(type) {
	case Raw:
		return string(c)
	case color.Color:
		return c.Code()
	case color.RGBColor:
		return autoConvertRGB(c)
	case color.Color256:
		return autoConvert256(c)
	default:
		panic(fmt.Errorf("unsupported color type: %T", c))
	}
}

func autoConvertRGB(rgb color.RGBColor) string {
	if color.SupportTrueColor() {
		return rgb.FullCode()
	}
	if color.Support256Color() {
		c256 := color.RgbTo256(rgb[0], rgb[1], rgb[2])
		return color.C256(c256, rgb[3] == 1).FullCode()
	}
	if color.SupportColor() {
		ansi := color.Rgb2basic(rgb[0], rgb[1], rgb[2], rgb[3] == 1)
		return color.Color(ansi).Code()
	}
	return ""
}

func autoConvert256(c256 color.Color256) string {
	if color.Support256Color() {
		return c256.FullCode()
	}
	if color.SupportColor() {
		// [color.Color256] does have a .Basic() function,
		// but its implementation does nothing and has a comment with "// TODO"
		rgb := c256.RGB()
		ansi := color.Rgb2basic(rgb[0], rgb[1], rgb[2], rgb[3] == 1)
		return color.Color(ansi).Code()
	}
	return ""
}
