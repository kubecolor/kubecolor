package color

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: replace with gookit/color
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
		return "light blue"
	case LightMagenta:
		return "light magenta"
	case LightCyan:
		return "light cyan"
	case LightWhite:
		return "light white"
	default:
		return fmt.Sprintf("%[1]T(%[1]d)", c)
	}
}

var colorParseReplacer = strings.NewReplacer(
	` `, "",
	`-`, "",
	`_`, "",
)

func Parse(s string) (Color, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return 0, nil
	}
	if num, err := strconv.Atoi(trimmed); err == nil {
		return Color(num), nil
	}
	// TODO: Replace this before merging with gookit/color parsing,
	// so we don't break available colors
	switch colorParseReplacer.Replace(strings.ToLower(trimmed)) {
	case "black":
		return Black, nil
	case "red":
		return Red, nil
	case "green":
		return Green, nil
	case "yellow", "orange":
		return Yellow, nil
	case "blue":
		return Blue, nil
	case "magenta", "purple":
		return Magenta, nil
	case "cyan":
		return Cyan, nil
	case "white":
		return White, nil
	case "lightblack":
		return LightBlack, nil
	case "lightred":
		return LightRed, nil
	case "lightgreen", "lime":
		return LightGreen, nil
	case "lightyellow", "lightorange":
		return LightYellow, nil
	case "lightblue":
		return LightBlue, nil
	case "lightmagenta", "lightpurple":
		return LightMagenta, nil
	case "lightcyan":
		return LightCyan, nil
	case "lightwhite":
		return LightWhite, nil
	default:
		return 0, fmt.Errorf("unknown color: %q", s)
	}
}

func (c Color) sequence() int {
	return int(c)
}

func Apply(val string, c Color) string {
	if c == 0 {
		return val
	}
	return fmt.Sprintf("%s[%dm%s%s[0m", escape, c.sequence(), val, escape)
}
