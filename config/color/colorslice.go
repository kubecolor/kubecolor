package color

import (
	"encoding"
	"fmt"
	"strings"
)

type ColorSlice []Color

var (
	_ encoding.TextMarshaler   = ColorSlice{}
	_ encoding.TextUnmarshaler = &ColorSlice{}
)

func MustParseColorSlice(s string) ColorSlice {
	c, err := ParseColorSlice(s)
	if err != nil {
		panic(fmt.Errorf("parse color slice: %w", err))
	}
	return c
}

func ParseColorSlice(s string) (ColorSlice, error) {
	split := strings.Split(s, "/")
	slice := make(ColorSlice, len(split))
	for i, sub := range split {
		col, err := ParseColor(strings.TrimSpace(sub))
		if err != nil {
			return nil, err
		}
		slice[i] = col
	}
	return slice, nil
}

func (s ColorSlice) String() string {
	strs := make([]string, len(s))
	for i, c := range s {
		strs[i] = c.String()
	}
	return strings.Join(strs, " / ")
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (s *ColorSlice) UnmarshalText(text []byte) error {
	newSlice, err := ParseColorSlice(string(text))
	if err != nil {
		return err
	}
	*s = newSlice
	return nil
}

// MarshalText implements [encoding.TextMarshaler].
func (s ColorSlice) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

func (s ColorSlice) ComputeCache() {
	for i := range s {
		s[i].ComputeCache()
	}
}
