package config

import (
	"encoding"
	"fmt"
	"strings"
)

type Paging string

const (
	// NOTE: When adding paging modes, remember to add them to [AllPagingModes] slice too.

	PagingAuto  Paging = "auto"
	PagingNever Paging = "never"
)

var (
	PagingDefault = PagingAuto

	AllPagingModes []Paging = []Paging{
		PagingAuto,
		PagingNever,
	}

	_ encoding.TextMarshaler   = PagingDefault
	_ encoding.TextUnmarshaler = &PagingDefault
)

func (p Paging) String() string {
	if p == "" {
		return "auto"
	}
	return string(p)
}

func ParsePaging(s string) (Paging, error) {
	if s == "" {
		return PagingAuto, nil
	}
	maybeValidPaging := Paging(strings.ToLower(s))
	for _, p := range AllPagingModes {
		if maybeValidPaging == p {
			return p, nil // reuse the interned string
		}
	}
	return PagingAuto, fmt.Errorf("invalid paging mode: %q", s)
}

// MarshalText implements [encoding.TextMarshaler].
func (p Paging) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (p *Paging) UnmarshalText(text []byte) error {
	newPaging, err := ParsePaging(string(text))
	if err != nil {
		return err
	}
	*p = newPaging
	return nil
}