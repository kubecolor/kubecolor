package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/xo/terminfo"
)

// ColorLevel is the color support for the terminal.
//
// This type is practically just wrapping [terminfo.ColorLevel], but allows
// us to have custom parsing and naming on the fields.
//
// It is defined here in "command" package instead of "config" package,
// as this setting shouldn't be available via the ~/.kube/color.yaml file.
type ColorLevel string

// ColorLevel values.
const (
	ColorLevelUnset     ColorLevel = ""
	ColorLevelNone      ColorLevel = "none"
	ColorLevelAuto      ColorLevel = "auto"
	ColorLevelBasic     ColorLevel = "basic"
	ColorLevel256       ColorLevel = "256"
	ColorLevelTrueColor ColorLevel = "truecolor"
)

func (c ColorLevel) TerminfoColorLevel() terminfo.ColorLevel {
	switch c {
	case ColorLevelNone:
		return terminfo.ColorLevelNone
	case ColorLevelBasic:
		return terminfo.ColorLevelBasic
	case ColorLevel256:
		return terminfo.ColorLevelHundreds
	case ColorLevelTrueColor:
		return terminfo.ColorLevelMillions

	case ColorLevelAuto:
		return color.TermColorLevel()
	default:
		return color.TermColorLevel()
	}
}

func (c ColorLevel) String() string {
	return string(c)
}

// UnmarshalText implements [encoding.TextUnmarshaller]
func (c *ColorLevel) UnmarshalText(text []byte) error {
	parsed, ok, err := ParseColorLevel(string(text))
	if err != nil {
		return err
	}
	if ok {
		*c = parsed
	}
	return nil
}

func ParseColorLevel(s string) (ColorLevel, bool, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "":
		return ColorLevelUnset, false, nil

	// Backward compatiblity, as the --force-colors flag was previously a bool flag
	case "true":
		return ColorLevelAuto, true, nil
	case "false":
		return ColorLevelUnset, false, nil

	case "none":
		return ColorLevelNone, true, nil
	case "auto":
		return ColorLevelAuto, true, nil

	// Parsing more than one valid input here, just to try be kinder to the user
	// in case they forgot the valid names
	case "basic", "3-bit", "3bit", "4-bit", "4bit":
		return ColorLevelBasic, true, nil
	case "256", "8-bit", "8bit":
		return ColorLevel256, true, nil
	case "truecolor", "true-color", "24-bit", "24bit":
		return ColorLevelTrueColor, true, nil
	default:
		return ColorLevelNone, false, fmt.Errorf("invalid color level: %q, must be one of: none, auto, basic, 256, truecolor", s)
	}
}

func parseColorLevelEnv(env string) (result ColorLevel, ok bool, err error) {
	result, ok, err = ParseColorLevel(os.Getenv(env))
	if err != nil {
		return ColorLevelNone, false, fmt.Errorf("parse env %s: %w", env, err)
	}
	return result, ok, err
}
