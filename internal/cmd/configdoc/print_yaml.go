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

type YAMLPrinter struct {
	*Program
}

func (p *YAMLPrinter) Print() error {
	fmt.Println("theme:")
	return p.printCategory(p.categories[0], []string{"theme"})
}

func (p *YAMLPrinter) printCategory(category Category, path []string) error {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, field := range category.Fields {
		newPath := append(path, field.Name)
		switch field.Type {
		case "color.Color":
			p.printField(tw, field, "color", newPath)
		case "color.Slice":
			p.printField(tw, field, "color[]", newPath)
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

func (p *YAMLPrinter) printField(w io.Writer, field Field, typeString string, path []string) error {
	fallback := p.formatFallback(field)
	desc := strings.ReplaceAll(field.Comment, "\n", " ")
	if fallback != "" {
		if desc != "" {
			desc += " "
		}
		desc += fallback
	}
	if desc != "" {
		desc = " " + desc
	}

	fmt.Fprintf(w, "%s%s: %s\t# (%s)%s\n",
		strings.Repeat("  ", len(path)-1),
		p.formatName(field),
		p.viper.GetString(strings.Join(path, ".")),
		typeString,
		desc,
	)

	return nil
}

func (p *YAMLPrinter) formatFallback(field Field) string {
	if field.DefaultFrom != "" {
		return fmt.Sprintf("(fallback to %s)", field.DefaultFrom)
	}
	if field.DefaultFromMany != "" {
		split := strings.Split(field.DefaultFromMany, ",")

		return fmt.Sprintf("(fallback to [%s])", strings.Join(split, " / "))
	}
	return ""
}

func (p *YAMLPrinter) formatName(field Field) string {
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
