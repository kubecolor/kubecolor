package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/invopop/jsonschema"
	"github.com/kubecolor/kubecolor/config"
)

var flags = struct {
	repo string
	out  string
}{
	repo: ".",
	out:  "-",
}

func init() {
	flag.StringVar(&flags.repo, "repo", flags.repo, "Path to root of code repository")
	flag.StringVar(&flags.out, "out", flags.out, "Where to write output. A single dash means stdout")
}

func main() {
	flag.Parse()

	r := &jsonschema.Reflector{
		Lookup:   Lookup,
		KeyNamer: Namer,
		Namer: func(t reflect.Type) string {
			return Namer(t.Name())
		},
		RequiredFromJSONSchemaTags: true,
		ExpandedStruct:             true,
	}

	r.AddGoComments("github.com/kubecolor/kubecolor", flags.repo)

	s := r.Reflect(&config.Config{})
	s.ID = "https://github.com/kubecolor/kubecolor/raw/main/config-schema.json"

	s.Definitions["color"] = &jsonschema.Schema{
		Type:        "string",
		Title:       "Color",
		Description: "A single color style, optionally setting foreground (text) color, background color, and/or modifier such as 'bold'.",
		Default:     "none",
		Examples: []any{
			"none",
			"red",
			"green",
			"yellow",
			"blue",
			"magenta",
			"cyan",
			"white",
			"black",
			"240",
			"aaff00",
			"#aaff00",
			"rgb(192, 255, 238)",
			"raw(4;53)",
			"gray:italic",
			"fg=white:bold:underline",
			"fg=yellow:bg=red:bold",
		},
	}

	s.Definitions["colorSlice"] = &jsonschema.Schema{
		Type:        "string",
		Title:       "Multiple colors",
		Description: "Allows multiple separate colors to be applied, separated by slash.",
		Examples: []any{
			"red/green/blue",
			"bg=red:underline/bg=green:italic/bg=blue:bold",
		},
	}

	s.Definitions["preset"] = &jsonschema.Schema{
		Type:        "string",
		Title:       "Color theme preset",
		Description: "Preset is a set of defaults for the color theme.",
		Default:     config.PresetDefault.String(),
		Enum:        castToAnySlice(config.AllPresets),
	}

	s.Definitions["paging"] = &jsonschema.Schema{
		Type:        "string",
		Title:       "Paging mode preference",
		Description: "Whether to pipe supported subcommands to a pager (\"auto\" or \"never\")",
		Default:     string(config.PagingDefault),
		Enum:        castToAnySlice(config.AllPagingModes),
	}

	s.Definitions["duration"] = &jsonschema.Schema{
		Type:        "string",
		Title:       "Time duration",
		Description: "A string value representing a time span, formatted as a Go time duration.",
		Default:     "0",
		Examples: []any{
			"30s",
			"5m",
			"10m",
			"1h30m",
			"5h",
		},
	}

	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	var out io.Writer = os.Stdout
	if flags.out != "-" {
		f, err := os.Create(flags.out)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		out = f
	}
	fmt.Fprintln(out, string(b))
	if flags.out != "-" {
		log.Println("Wrote to:", flags.out)
	}
}

func castToAnySlice[E any](s []E) []any {
	slice := make([]any, len(s))
	for i, v := range s {
		slice[i] = v
	}
	return slice
}

// Lookup allows a function to be defined that will provide a custom mapping of
// types to Schema IDs.
func Lookup(t reflect.Type) jsonschema.ID {
	switch t.Name() {
	case "Color", "Slice", "Preset", "Paging", "Duration":
		return jsonschema.ID("#/$defs/" + Namer(t.Name()))
	default:
		return ""
	}
}

// Namer allows customizing of type names.
func Namer(s string) string {
	switch s {
	case "GUID":
		return "guid"
	case "Slice":
		return "colorSlice"
	}
	var sb strings.Builder
	sb.Grow(len(s))
	firstRune, size := utf8.DecodeRuneInString(s)
	sb.WriteRune(unicode.ToLower(firstRune))
	sb.WriteString(s[size:])
	return sb.String()
}
