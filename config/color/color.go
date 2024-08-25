package color

import (
	"cmp"
	"encoding"
	"fmt"
	"regexp"
	"strings"

	"github.com/gookit/color"
)

type ColorCode interface {
	Code() string
}

type Color struct {
	Source string
	Parsed []ColorCode

	cached     bool
	cachedCode string
}

var (
	_ encoding.TextMarshaler   = Color{}
	_ encoding.TextUnmarshaler = &Color{}
)

// String returns the display name of this color.
func (c Color) String() string {
	if c.cached && c.cachedCode == "" {
		return "none"
	}
	return cmp.Or(c.Source, c.cachedCode)
}

// IsNoop returns true when no coloring would be applied by this Color.
func (c Color) IsNoop() bool {
	if !c.cached {
		c.ComputeCache()
	}
	return c.cachedCode == ""
}

// ANSICode returns the ASNI coloring code that this color would apply.
func (c Color) ANSICode() string {
	if !c.cached {
		c.ComputeCache()
	}
	return c.cachedCode
}

// Render returns the string wrapped in color codes from this color.
func (c Color) Render(s string) string {
	if !c.cached {
		c.ComputeCache()
	}
	if strings.ContainsRune(s, '\033') {
		return c.renderInject(s)
	}
	return color.RenderString(c.cachedCode, s)
}

var afterResetRegex = regexp.MustCompile("\033\\[0m[^\033]")

func (c Color) renderInject(s string) string {
	if strings.HasPrefix(s, "\033[") &&
		strings.HasSuffix(s, "\033[0m") &&
		strings.Count(s, "\033[") == 2 {
		// If full string is colored, then doesn't matter if we add colors
		return s
	}
	updated := afterResetRegex.ReplaceAllStringFunc(s, func(orig string) string {
		lastByte := orig[len(orig)-1]
		return fmt.Sprintf("\033[0m\033[%sm%c", c.cachedCode, lastByte)
	})
	return color.RenderString(c.cachedCode, updated)
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
	colors := make([]ColorCode, 0, count)

	// Using ':' separator to mimic Linux LS_COLORS
	fields := strings.FieldsFunc(s, func(r rune) bool { return r == ':' })
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		if field == "none" {
			continue
		}

		c, err := parseColorField(field)
		if err != nil {
			return Color{}, err
		}
		if c == nil {
			continue
		}

		colors = append(colors, c)
	}

	return Color{
		Source: s,
		Parsed: colors,
	}, nil
}

// Set implements [flag.Value].
func (c *Color) Set(text string) error {
	newColor, err := ParseColor(text)
	if err != nil {
		return fmt.Errorf("parse color: %w", err)
	}
	*c = newColor
	return nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (c *Color) UnmarshalText(text []byte) error {
	return c.Set(string(text))
}

// MarshalText implements [encoding.TextMarshaler].
func (c Color) MarshalText() (text []byte, err error) {
	return []byte(c.String()), nil
}

func (c Color) IsZero() bool {
	return c.Source == ""
}

func (c *Color) ComputeCache() {
	codes := make([]string, len(c.Parsed))
	for i, parsed := range c.Parsed {
		codes[i] = ConvertColorCode(parsed)
	}
	c.cachedCode = strings.Join(codes, ";")
	c.cached = true
}

func ClearCode(s string) string {
	return color.ClearCode(s)
}

func ForceColor() {
	color.ForceColor()
	color.Enable = true
}
