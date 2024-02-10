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
)

const (
	LightBlack Color = iota + 90 // Grey...
	LightRed
	LightGreen
	LightYellow
	LightBlue
	LightMagenta
	LightCyan
	LightWhite // White...
)

func (c Color) String() string {
	switch c {
	case Black:
		return "black"
	case Red:
		return "red"
	case Green:
		return "green"
	case Yellow:
		return "yellow"
	case Blue:
		return "blue"
	case Magenta:
		return "magenta"
	case Cyan:
		return "cyan"
	case White:
		return "white"
	case LightBlack:
		return "light black"
	case LightRed:
		return "light red"
	case LightGreen:
		return "light green"
	case LightYellow:
		return "light yellow"
	case LightBlue:
		return "blue"
	case LightMagenta:
		return "magenta"
	case LightCyan:
		return "cyan"
	case LightWhite:
		return "white"
	default:
		return fmt.Sprintf("%[1]T(%[1]d)", c)
	}
}

func (c Color) sequence() int {
	return int(c)
}

func Apply(val string, c Color) string {
	return fmt.Sprintf("%s[%dm%s%s[0m", escape, c.sequence(), val, escape)
}
