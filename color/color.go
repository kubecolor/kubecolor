package color

import (
	"fmt"
)

type Color int

const escape = "\x1b"

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White

	LightBlack Color = iota + 90 // Grey...
	LightRed
	LightGreen
	LightYellow
	LightBlue
	LightMagenta
	LightCyan
	LightWhite // White...
)

var (
	basicColorMap = map[string]Color{
		"black":        Black,
		"red":          Red,
		"green":        Green,
		"yellow":       Yellow,
		"blue":         Blue,
		"magenta":      Magenta,
		"cyan":         Cyan,
		"white":        White,
		"lightblack":   LightBlack,
		"lightred":     LightRed,
		"lightgreen":   LightGreen,
		"lightyellow":  LightYellow,
		"lightblue":    LightBlue,
		"lightmagenta": LightMagenta,
		"lightcyan":    LightCyan,
		"lightwhite":   LightWhite,
	}
)

func (c Color) sequence() int {
	return int(c)
}

func Apply(val string, c Color) string {
	return fmt.Sprintf("%s[%dm%s%s[0m", escape, c.sequence(), val, escape)
}

// StringToColor returns the color int from the specified color name
func StringToColor(s string) Color {
	return basicColorMap[s]
}
