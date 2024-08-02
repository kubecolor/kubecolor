package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/kubecolor/kubecolor/internal/stringutil"
)

func parseColorField(field string) (ColorCode, error) {
	key, value, hasEqualSign := strings.Cut(field, "=")
	key = strings.TrimSpace(key)

	if !hasEqualSign {
		if c := parseColorOp(key); c != 0 {
			return c, nil
		}
		return parseColorFg(key)
	}

	value = strings.TrimSpace(value)
	switch key {
	case "fg":
		return parseColorFg(value)
	case "bg":
		return parseColorBg(value)
	default:
		return nil, fmt.Errorf("invalid color key %q", key+"=")
	}
}

func parseColorOp(s string) color.Color {
	switch strings.ToLower(s) {
	case "reset":
		return color.OpReset
	case "bold", "b":
		return color.Bold
	case "fuzzy":
		return color.OpFuzzy
	case "italic", "i":
		return color.OpItalic
	case "underscore", "underscored", "u", "underline", "underlined":
		return color.OpUnderscore
	case "blink":
		return color.OpBlink
	case "fastblink":
		return color.OpFastBlink
	case "reverse", "invert", "inverted":
		return color.OpReverse
	case "concealed", "hidden", "invisible":
		return color.OpConcealed
	case "strikethrough", "strike":
		return color.OpStrikethrough
	default:
		return 0
	}
}

var colorNameReplacer = strings.NewReplacer(
	"hi_", "hi",
	"hi-", "hi",
	"light_", "light",
	"light-", "light",
)

func parseColorFg(s string) (ColorCode, error) {
	if c := parseColorFgName(s); c != 0 {
		return c, nil
	}
	coder, err := parseColorSyntax(s, false)
	if err != nil {
		return nil, fmt.Errorf("bg: %w", err)
	}
	return coder, nil
}

func parseColorFgName(s string) color.Color {
	switch strings.ToLower(colorNameReplacer.Replace(s)) {
	// Basic colors
	case "black":
		return color.FgBlack
	case "red":
		return color.FgRed
	case "green":
		return color.FgGreen
	case "brown", "yellow":
		return color.FgYellow
	case "blue":
		return color.FgBlue
	case "magenta", "purple":
		return color.FgMagenta
	case "cyan":
		return color.FgCyan
	case "white":
		return color.FgWhite
	case "default", "normal":
		return color.FgDefault // no color

	// Light/high colors
	case "hiblack", "lightblack",
		"darkgray", "gray",
		"darkgrey", "grey":
		return color.FgGray
	case "hired", "lightred":
		return color.FgLightRed
	case "higreen", "lightgreen", "lime":
		return color.FgLightGreen
	case "hibrown", "lightbrown",
		"hiyellow", "lightyellow",
		"gold":
		return color.FgLightYellow
	case "hiblue", "lightblue":
		return color.FgLightBlue
	case "himagenta", "lightmagenta",
		"hipurple", "lightpurple":
		return color.FgLightMagenta
	case "hicyan", "lightcyan":
		return color.FgLightCyan
	case "hiwhite", "lightwhite":
		return color.FgLightWhite
	default:
		return 0
	}
}

func parseColorBg(s string) (ColorCode, error) {
	if c := parseColorBgName(s); c != 0 {
		return c, nil
	}
	coder, err := parseColorSyntax(s, true)
	if err != nil {
		return nil, fmt.Errorf("bg: %w", err)
	}
	return coder, nil
}

func parseColorBgName(s string) color.Color {
	switch strings.ToLower(colorNameReplacer.Replace(s)) {
	// Basic colors
	case "black":
		return color.BgBlack
	case "red":
		return color.BgRed
	case "green":
		return color.BgGreen
	case "brown", "yellow":
		return color.BgYellow
	case "blue":
		return color.BgBlue
	case "magenta", "purple":
		return color.BgMagenta
	case "cyan":
		return color.BgCyan
	case "white":
		return color.BgWhite
	case "default", "normal", "transparent":
		return color.BgDefault // no color

	// Light/high colors
	case "hiblack", "lightblack",
		"darkgray", "gray",
		"darkgrey", "grey":
		return color.BgGray
	case "hired", "lightred":
		return color.BgLightRed
	case "higreen", "lightgreen", "lime":
		return color.BgLightGreen
	case "hibrown", "lightbrown",
		"hiyellow", "lightyellow",
		"gold":
		return color.BgLightYellow
	case "hiblue", "lightblue":
		return color.BgLightBlue
	case "himagenta", "lightmagenta",
		"hipurple", "lightpurple":
		return color.BgLightMagenta
	case "hicyan", "lightcyan":
		return color.BgLightCyan
	case "hiwhite", "lightwhite":
		return color.BgLightWhite
	default:
		return 0
	}
}

var (
	isHexRegex     = regexp.MustCompile(`^(#|0x)?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)
	is256CodeRegex = regexp.MustCompile(`^[0-9]{1,3}$`)
)

// parseColorSyntax is derived from [github.com/gookit/color] package:
// [https://github.com/gookit/color/blob/9027b9d2a5168ea482a8a8b46711191450514aa3/color_tag.go#L453-L471]
func parseColorSyntax(s string, isBg bool) (ColorCode, error) {
	rgb, ok, err := parseColorFunctionRGB(s, isBg)
	if err != nil {
		return nil, fmt.Errorf("parse rgb(...): %w", err)
	}
	if ok {
		return rgb, nil
	}

	raw, ok, err := parseColorFunction(s, "raw")
	if err != nil {
		return nil, fmt.Errorf("parse raw(...): %w", err)
	}
	if ok {
		return Raw(raw), nil
	}

	if strings.Count(s, ",") == 2 {
		// parse again as rgb, but treat "255, 200, 100" as "rgb(255, 200, 100)"
		// but discard errors, as this is not a precise syntax
		if rgb, ok, err := parseColorFunctionRGB(fmt.Sprintf("rgb(%s)", s), isBg); err == nil && ok {
			return rgb, nil
		}
	}

	if len(s) < 4 && is256CodeRegex.MatchString(s) { // single 256 code
		num, err := strconv.ParseUint(s, 10, 8)
		if err == nil {
			c := color.C256(uint8(num), isBg)
			return c, nil
		}
	}

	if isHexRegex.MatchString(s) { // hex: "#fc1cac"
		hex := color.HEX(s, isBg)
		return hex, nil
	}

	return nil, fmt.Errorf("invalid color format: %q", s)
}

func parseColorFunctionRGB(s string, isBg bool) (color.RGBColor, bool, error) {
	rgbStr, ok, err := parseColorFunction(s, "rgb")
	if err != nil || !ok {
		return color.RGBColor{}, false, err
	}
	rgb, err := parse3Uints(rgbStr, 8)
	if err != nil {
		return color.RGBColor{}, true, err
	}
	c := color.RGB(uint8(rgb[0]), uint8(rgb[1]), uint8(rgb[2]), isBg)
	return c, true, nil
}

// parseColorFunction parses a function syntax, such as:
//
//	rgb(255, 222, 100) => "255, 222, 100"
//	hsl(0.5, 0.5, 0.5) => "0.5, 0.5, 0.5"
func parseColorFunction(s, name string) (string, bool, error) {
	withoutName, ok := strings.CutPrefix(s, name)
	if !ok {
		return "", false, nil
	}
	withoutStart, ok := strings.CutPrefix(withoutName, "(")
	if !ok {
		return "", false, fmt.Errorf(`missing opening parentheses "(": got %q`, s)
	}
	onlyValues, ok := strings.CutSuffix(withoutStart, ")")
	if !ok {
		return "", false, fmt.Errorf(`missing closing parentheses ")": got %q`, s)
	}
	if strings.TrimSpace(onlyValues) == "" {
		return "", false, fmt.Errorf("color function must not be empty")
	}
	return onlyValues, true, nil
}

func parse3Uints(s string, bitSize int) ([3]uint64, error) {
	split := stringutil.SplitAndTrimSpace(s, ",")
	if len(split) != 3 {
		return [3]uint64{}, fmt.Errorf(`invalid number of args; want 3, got: %d`, len(split))
	}
	i1, err := strconv.ParseUint(split[0], 10, bitSize)
	if err != nil {
		return [3]uint64{}, fmt.Errorf("arg 1/3: %w", err)
	}
	i2, err := strconv.ParseUint(split[1], 10, bitSize)
	if err != nil {
		return [3]uint64{}, fmt.Errorf("arg 2/3: %w", err)
	}
	i3, err := strconv.ParseUint(split[2], 10, bitSize)
	if err != nil {
		return [3]uint64{}, fmt.Errorf("arg 3/3: %w", err)
	}
	return [3]uint64{i1, i2, i3}, nil
}
