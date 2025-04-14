package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"unicode"
	"unicode/utf8"
)

type FilePrinter struct {
	*Program
}

func (p *FilePrinter) Print() error {
	fmt.Println("theme:")
	return p.printCategory(p.categories[0], []string{"theme"})
}

func (p *FilePrinter) printCategory(category Category, path []string) error {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, field := range category.Fields {
		newPath := append(path, field.Name)
		switch field.Type {
		case "color.Color":
			p.printField(tw, field, "color", newPath)
		case "color.Slice":
			p.printField(tw, field, "color-list", newPath)
		default:
			sub, ok := p.findCategory(field.Type)
			if !ok {
				return fmt.Errorf("invalid category field type: %q", field.Type)
			}
			fmt.Fprintf(tw, "%s%s:\n",
				strings.Repeat("  ", len(path)),
				p.formatName(field),
			)
			tw.Flush()
			if err := p.printCategory(sub, newPath); err != nil {
				return err
			}
		}
	}
	tw.Flush()
	return nil
}

func (p *FilePrinter) printField(w io.Writer, field Field, typeString string, path []string) error {
	fallback := p.formatFallback(field)
	desc := strings.ReplaceAll(field.Comment, "\n", " ")

	indent := strings.Repeat("  ", len(path)-1)

	if desc != "" {
		fmt.Fprintf(w, "%s# %s\n",
			indent,
			desc)
	}
	if fallback == "" {
		fmt.Fprintf(w, "%s%s: !%s %s\n",
			indent,
			p.formatName(field),
			typeString,
			p.viper.GetString(strings.Join(path, ".")))
	} else {
		fmt.Fprintf(w, "%s%s: !%s\t# default = %s\n",
			indent,
			p.formatName(field),
			typeString,
			fallback)
	}
	fmt.Fprintln(w)

	return nil
}

func (p *FilePrinter) formatFallback(field Field) string {
	if field.DefaultFrom != "" {
		return fmt.Sprintf("$%s", field.DefaultFrom)
	}
	if field.DefaultFromMany != "" {
		split := strings.Split(field.DefaultFromMany, ",")
		for i, v := range split {
			split[i] = "$" + v
		}
		return fmt.Sprintf("[%s])", strings.Join(split, " / "))
	}
	return ""
}

func (p *FilePrinter) formatName(field Field) string {
	lower := camelCase(field.Name)

	switch lower {
	case "true", "false", "null":
		return fmt.Sprintf(`"%s"`, lower)
	default:
		return lower
	}
}

func camelCase(s string) string {
	if isAllUppercase(s) {
		return strings.ToLower(s)
	}
	return firstLetterLowercase(s)
}

func firstLetterLowercase(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	firstRune, size := utf8.DecodeRuneInString(s)
	sb.WriteRune(unicode.ToLower(firstRune))
	sb.WriteString(s[size:])
	return sb.String()
}

func isAllUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return false
		}
	}
	return true
}
