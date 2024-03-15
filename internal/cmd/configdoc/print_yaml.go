package main

import (
	"fmt"
	"strings"
)

type YAMLPrinter struct {
	*Program
}

func (p *YAMLPrinter) Print() error {
	fmt.Println("theme:")
	p.printCategory(p.categories[0], []string{"theme"})
	return nil
}

func (p *YAMLPrinter) printCategory(category Category, path []string) error {
	for _, field := range category.Fields {
		newPath := append(path, field.Name)
		switch field.Type {
		case "Color":
			p.printField(field, "color", newPath)
		case "ColorSlice":
			p.printField(field, "color[]", newPath)
		default:
			sub, ok := p.findCategory(field.Type)
			if !ok {
				return fmt.Errorf("invalid category field type: %q", field.Type)
			}
			fmt.Printf("%s%s:\n",
				strings.Repeat("  ", len(path)),
				strings.ToLower(field.Name),
			)
			if err := p.printCategory(sub, newPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *YAMLPrinter) printField(field Field, typeString string, path []string) error {

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

	fmt.Printf("%s%s: %s # (%s)%s\n",
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
	lower := strings.ToLower(field.Name)
	switch lower {
	case "true", "false", "null":
		return fmt.Sprintf(`"%s"`, lower)
	default:
		return lower
	}
}
