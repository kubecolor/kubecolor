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

func (c Color) sequence() int {
	return int(c)
}

func Apply(val string, c Color) string {
	return fmt.Sprintf("%s[%dm%s%s[0m", escape, c.sequence(), val, escape)
}
