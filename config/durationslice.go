package config

import (
	"encoding"
	"fmt"
	"strings"
	"time"

	"github.com/kubecolor/kubecolor/internal/stringutil"
)

// DurationSlice is an ordered list of time durations, parsed from a
// "/"-separated string such as "5m/2h/1d". It mirrors [config/color.Slice]
// so the two lists (age thresholds and their colors) share the same syntax.
type DurationSlice []time.Duration

var (
	_ encoding.TextMarshaler   = DurationSlice{}
	_ encoding.TextUnmarshaler = &DurationSlice{}
)

// ParseDurationSlice parses a "/"-separated list of durations. Each element is
// parsed as a human-friendly k8s duration ("5m", "2h", "1d", "7d", "1y"),
// falling back to Go's [time.ParseDuration] for values it doesn't handle such
// as sub-second units ("500ms"). An empty string yields an empty slice.
func ParseDurationSlice(s string) (DurationSlice, error) {
	if strings.TrimSpace(s) == "" {
		return DurationSlice{}, nil
	}
	split := strings.Split(s, "/")
	slice := make(DurationSlice, len(split))
	for i, sub := range split {
		d, err := parseDuration(strings.TrimSpace(sub))
		if err != nil {
			return nil, err
		}
		slice[i] = d
	}
	return slice, nil
}

func parseDuration(s string) (time.Duration, error) {
	if d, ok := stringutil.ParseHumanDuration(s); ok {
		return d, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("parse duration %q: %w", s, err)
	}
	return d, nil
}

func MustParseDurationSlice(s string) DurationSlice {
	slice, err := ParseDurationSlice(s)
	if err != nil {
		panic(fmt.Errorf("parse duration slice: %w", err))
	}
	return slice
}

func (s DurationSlice) String() string {
	strs := make([]string, len(s))
	for i, d := range s {
		strs[i] = d.String()
	}
	return strings.Join(strs, " / ")
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (s *DurationSlice) UnmarshalText(text []byte) error {
	slice, err := ParseDurationSlice(string(text))
	if err != nil {
		return err
	}
	*s = slice
	return nil
}

// MarshalText implements [encoding.TextMarshaler].
func (s DurationSlice) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}
