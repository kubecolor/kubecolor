package color

import (
	"encoding"
	"fmt"
	"strings"
)

type Slice []Color

var (
	_ encoding.TextMarshaler   = Slice{}
	_ encoding.TextUnmarshaler = &Slice{}
)

func MustParseSlice(s string) Slice {
	c, err := ParseSlice(s)
	if err != nil {
		panic(fmt.Errorf("parse color slice: %w", err))
	}
	return c
}

func ParseSlice(s string) (Slice, error) {
	split := strings.Split(s, "/")
	slice := make(Slice, len(split))
	for i, sub := range split {
		col, err := Parse(strings.TrimSpace(sub))
		if err != nil {
			return nil, err
		}
		slice[i] = col
	}
	return slice, nil
}

func (s Slice) String() string {
	strs := make([]string, len(s))
	for i, c := range s {
		strs[i] = c.String()
	}
	return strings.Join(strs, " / ")
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (s *Slice) UnmarshalText(text []byte) error {
	newSlice, err := ParseSlice(string(text))
	if err != nil {
		return err
	}
	*s = newSlice
	return nil
}

// MarshalText implements [encoding.TextMarshaler].
func (s Slice) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

func (s Slice) ComputeCache() {
	for i := range s {
		s[i].ComputeCache()
	}
}
