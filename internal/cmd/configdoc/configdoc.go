package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/kubecolor/kubecolor/config"
	"github.com/spf13/viper"
)

var flags = struct {
	file string
}{
	file: filepath.Join("config", "theme.go"),
}

func init() {
	flag.StringVar(&flags.file, "file", flags.file, "Path to theme.go file")
}

func main() {
	flag.Parse()

	prog := Program{}
	if err := prog.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

type Program struct {
	viper      *viper.Viper
	categories []Category

	writer            *tabwriter.Writer
	lastPrintWasValue bool
}

func (p *Program) Run() error {
	v := config.NewViper()
	if err := config.ApplyThemePreset(v); err != nil {
		return err
	}
	p.viper = v

	categories, err := ParseCategories()
	if err != nil {
		return fmt.Errorf("parse categories: %w", err)
	}
	p.categories = categories
	if err := p.Print(); err != nil {
		return fmt.Errorf("print: %w", err)
	}
	return nil
}

func (p *Program) Print() error {
	p.writer = tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', 0)

	p.printColumns("Environment variable", "Type", "Description", "Dark theme")
	p.printColumns("--------------------", "----", "-----------", "----------")

	p.printCategory(p.categories[0], []string{"theme"})

	p.writer.Flush()
	return nil
}

func (p *Program) printCategory(category Category, path []string) error {
	for i, field := range category.Fields {
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
			if p.lastPrintWasValue && i != 0 {
				p.printColumns("", "", "", "")
			}
			if err := p.printCategory(sub, newPath); err != nil {
				return err
			}
			isLast := i == len(category.Fields)-1
			if p.lastPrintWasValue && !isLast {
				p.printColumns("", "", "", "")
			}
		}
	}
	return nil
}

func (p *Program) printField(field Field, typeString string, path []string) error {
	value := p.viper.GetString(strings.Join(path, "."))
	darkTheme := fmt.Sprintf("`%s`", value)
	if value == "" {
		darkTheme = ""
	}

	fallback := formatFallback(field)
	desc := strings.ReplaceAll(field.Comment, "\n", "<br/>")
	if fallback != "" {
		if desc != "" {
			desc += "<br/>"
		}
		desc += fallback
	}

	p.printColumns(fmt.Sprintf("`%s`", pathString(path)), typeString, desc, darkTheme)
	return nil
}

func formatFallback(field Field) string {
	if field.DefaultFrom != "" {
		defaultFrom := pathString(strings.Split(field.DefaultFrom, "."))
		return fmt.Sprintf("*(fallback to `%s`)*", defaultFrom)
	}
	if field.DefaultFromMany != "" {
		split := strings.Split(field.DefaultFromMany, ",")
		for i, s := range split {
			split[i] = pathString(strings.Split(s, "."))
		}

		return fmt.Sprintf("*(fallback to `[%s]`)*", strings.Join(split, " / "))
	}
	return ""
}

func (p *Program) printColumns(envVar, typ, desc, darkTheme string) {
	if darkTheme == "" {
		fmt.Fprintf(p.writer, "|\t%s\t|\t%s\t|\t%s\t|\n", envVar, typ, desc)
	} else {
		fmt.Fprintf(p.writer, "|\t%s\t|\t%s\t|\t%s\t|\t%s\n", envVar, typ, desc, darkTheme)
	}
	p.lastPrintWasValue = envVar != ""
}

func pathString(path []string) string {
	var sb strings.Builder
	sb.WriteString("KUBECOLOR")
	for _, s := range path {
		sb.WriteByte('_')
		sb.WriteString(strings.ToUpper(s))
	}
	return sb.String()
}
