package flagutil

import (
	"encoding"
	"fmt"
	"slices"
	"strings"
)

type FlagSet []*Flag

func (s *FlagSet) add(flag *Flag) *Flag {
	*s = append(*s, flag)
	return flag
}

func (s *FlagSet) NewString(name, desc string) *Flag {
	return s.add(&Flag{
		Name:        name,
		Description: desc,
	})
}

func (s *FlagSet) NewBool(name, desc string) *Flag {
	return s.add(&Flag{
		Name:        name,
		Description: desc,
		Enum:        []string{"true", "false"},
	})
}

func (s *FlagSet) NewUnmarshaller(name, desc string, value encoding.TextUnmarshaler) *Flag {
	return s.add(&Flag{
		Name:         name,
		Description:  desc,
		Unmarshaller: value,
	})
}

func (s *FlagSet) ParseArg(arg string) (*Flag, error) {
	if !strings.HasPrefix(arg, "--") {
		return nil, nil
	}
	name, value, split := strings.Cut(arg, "=")
	for _, f := range *s {
		if name != f.Name {
			continue
		}
		f.Value = value

		if split {
			switch {
			case len(f.Enum) > 0:
				if !slices.Contains(f.Enum, value) {
					return f, fmt.Errorf(`flag %s: must be one of: %s`, name, strings.Join(f.Enum, ", "))
				}

			case f.Unmarshaller != nil:
				if err := f.Unmarshaller.UnmarshalText([]byte(value)); err != nil {
					return f, fmt.Errorf(`flag %s: %w`, name, err)
				}
			}
		}

		return f, nil
	}
	return nil, nil
}

type Flag struct {
	Name         string
	Description  string
	Enum         []string
	Unmarshaller encoding.TextUnmarshaler

	Value string
}

func (f *Flag) BoolValue() bool {
	if f.Value == "false" {
		return false
	}
	// bool flags treat no value as true (e.g "--plain" is same as "--plain=true")
	return true
}
