package config

import (
	"cmp"
	"encoding"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gookit/color"
)

type Color struct {
	Source string
	Code   string
}

var (
	_ encoding.TextMarshaler   = Color{}
	_ encoding.TextUnmarshaler = &Color{}
)

// String returns the display name of this color.
func (c Color) String() string {
	return cmp.Or(c.Source, c.Code)
}

// Render returns the string wrapped in color codes from this color.
func (c Color) Render(s string) string {
	return color.RenderString(c.Code, s)
}

// Sprint returns the stringified args (concatenated right after each other)
// wrapped in color codes from this color.
func (c Color) Sprint(args ...any) string {
	return c.Render(fmt.Sprint(args...))
}

// Sprintln returns the stringified args (concatenated right after each other)
// wrapped in color codes from this color, as well as a trailing newline.
// The newline character is added after the color reset.
func (c Color) Sprintln(args ...any) string {
	// Don't use [fmt.Sprintln] here because we want the \e[0m to be before the \n
	return c.Render(fmt.Sprint(args...)) + "\n"
}

// Sprintf returns the formatted string wrapped in color codes from this color.
func (c Color) Sprintf(format string, args ...any) string {
	return c.Render(fmt.Sprintf(format, args...))
}

func MustParseColor(s string) Color {
	c, err := ParseColor(s)
	if err != nil {
		panic(fmt.Errorf("parse color: %w", err))
	}
	return c
}

func ParseColor(s string) (Color, error) {
	// Heavily inspired by color.ParseCodeFromAttr
	s = strings.Trim(s, ":,")

	count := strings.Count(s, ":")
	codes := make([]string, 0, count)

	// Using ':' separator to mimic Linux LS_COLORS
	fields := strings.FieldsFunc(s, func(r rune) bool { return r == ':' })
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		c, err := parseColorField(field)
		if err != nil {
			return Color{}, err
		}
		if c == "" {
			continue
		}

		codes = append(codes, c)
	}

	return Color{
		Source: s,
		Code:   strings.Join(codes, ";"),
	}, nil
}

func parseColorField(field string) (string, error) {
	key, value, hasEqualSign := strings.Cut(field, "=")
	key = strings.TrimSpace(key)

	if !hasEqualSign {
		if c := parseColorOp(key); c != 0 {
			return c.Code(), nil
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
		return "", fmt.Errorf("invalid color key %q", key+"=")
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
	case "reverse", "inverted":
		return color.OpReverse
	case "concealed", "hidden", "invisible":
		return color.OpConcealed
	case "strikethrough":
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

func parseColorFg(s string) (string, error) {
	if c := parseColorFgName(s); c != 0 {
		return c.Code(), nil
	}
	code, err := parseColorSyntax(s, false)
	if err != nil {
		return "", err
	}
	if code == "" {
		return "", fmt.Errorf("invalid fg color: %q", s)
	}
	return code, nil
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
	case "higreen", "lightgreen",
		"lime":
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

func parseColorBg(s string) (string, error) {
	if c := parseColorBgName(s); c != 0 {
		return c.Code(), nil
	}
	code, err := parseColorSyntax(s, true)
	if err != nil {
		return "", err
	}
	if code == "" {
		return "", fmt.Errorf("invalid bg color: %q", s)
	}
	return code, nil
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
	case "default", "normal", "transparent", "none":
		return color.BgDefault // no color

	// Light/high colors
	case "hiblack", "lightblack",
		"darkgray", "gray",
		"darkgrey", "grey":
		return color.BgGray
	case "hired", "lightred":
		return color.BgLightRed
	case "higreen", "lightgreen",
		"lime":
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
func parseColorSyntax(s string, isBg bool) (string, error) {
	if isHexRegex.MatchString(s) { // hex: "#fc1cac"
		return color.HEX(s, isBg).String(), nil
	}

	rgb, err := parseColorFunctionRGB(s, isBg)
	if err != nil {
		return "", fmt.Errorf("parse rgb(...): %w", err)
	}
	if rgb != "" {
		return rgb, nil
	}

	hsl, err := parseColorFunctionHSL(s, isBg)
	if err != nil {
		return "", fmt.Errorf("parse hsl(...): %w", err)
	}
	if hsl != "" {
		return hsl, nil
	}

	if strings.Count(s, ",") == 2 {
		// parse again as rgb, but treat "255, 200, 100" as "rgb(255, 200, 100)"
		// but discard errors, as this is not a precise syntax
		if rgb, err := parseColorFunctionRGB("rgb("+s+")", isBg); err == nil && rgb != "" {
			return rgb, nil
		}
	}

	if len(s) < 4 && is256CodeRegex.MatchString(s) { // single 256 code
		if isBg {
			return color.Bg256Pfx + s, nil
		} else {
			return color.Fg256Pfx + s, nil
		}
	}
	return "", nil
}

// parseColorFunction parses a function syntax, such as:
//
//	rgb(255, 222, 100)
//	hsl(0.5, 0.5, 0.5)
func parseColorFunction(s, name string) ([3]string, bool, error) {
	withoutName, ok := strings.CutPrefix(s, name)
	if !ok {
		return [3]string{}, false, nil
	}
	withoutStart, ok := strings.CutPrefix(withoutName, "(")
	if !ok {
		return [3]string{}, true, fmt.Errorf(`missing opening parentheses "(": got %q`, s)
	}
	onlyValues, ok := strings.CutSuffix(withoutStart, ")")
	if !ok {
		return [3]string{}, true, fmt.Errorf(`missing closing parentheses ")": got %q`, s)
	}
	split := strings.Split(onlyValues, ",")
	if len(split) != 3 {
		return [3]string{}, true, fmt.Errorf(`invalid number of args; want 3, got: %d`, len(split))
	}
	return [3]string{
		strings.TrimSpace(split[0]),
		strings.TrimSpace(split[1]),
		strings.TrimSpace(split[2]),
	}, true, nil
}

func parseColorFunctionHSL(s string, isBg bool) (string, error) {
	hsl, ok, err := parseColorFunction(s, "hsl")
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	hue, err := strconv.ParseFloat(hsl[0], 32)
	if err != nil {
		return "", fmt.Errorf("hsl.h: %w", err)
	}
	saturation, err := strconv.ParseFloat(hsl[1], 32)
	if err != nil {
		return "", fmt.Errorf("hsl.s: %w", err)
	}
	lightness, err := strconv.ParseFloat(hsl[2], 32)
	if err != nil {
		return "", fmt.Errorf("hsl.l: %w", err)
	}
	return color.HSL(hue, saturation, lightness, isBg).FullCode(), nil
}

func parseColorFunctionRGB(s string, isBg bool) (string, error) {
	rgb, ok, err := parseColorFunction(s, "rgb")
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	r, err := strconv.ParseUint(rgb[0], 10, 8)
	if err != nil {
		return "", fmt.Errorf("rgb.r: %w", err)
	}
	g, err := strconv.ParseUint(rgb[1], 10, 8)
	if err != nil {
		return "", fmt.Errorf("rgb.g: %w", err)
	}
	b, err := strconv.ParseUint(rgb[2], 10, 8)
	if err != nil {
		return "", fmt.Errorf("rgb.b: %w", err)
	}
	return color.RGB(uint8(r), uint8(g), uint8(b), isBg).FullCode(), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (c *Color) UnmarshalText(text []byte) error {
	newColor, err := ParseColor(string(text))
	if err != nil {
		return fmt.Errorf("parse color: %w", err)
	}
	*c = newColor
	return nil
}

// MarshalText implements [encoding.TextMarshaler].
func (c Color) MarshalText() (text []byte, err error) {
	return []byte(c.String()), nil
}
