package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
)

type EnvPrinter struct {
	*Program
	writer            *tabwriter.Writer
	lastPrintWasValue bool
}

func (p *EnvPrinter) Print() error {
	return p.printCategory(p.categories[0], []string{"theme"})
}

func (p *EnvPrinter) printCategory(category Category, path []string) error {
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
			fmt.Fprintln(tw)
			tw.Flush()
			if err := p.printCategory(sub, newPath); err != nil {
				return err
			}
		}
	}
	tw.Flush()
	return nil
}

func (p *EnvPrinter) printField(w io.Writer, field Field, typeString string, path []string) error {
	fallback := p.formatFallback(field)
	desc := strings.ReplaceAll(field.Comment, "\n", " ")

	if desc != "" {
		fmt.Fprintf(w, "# (%s) %s\n",
			typeString,
			desc)
	}
	if fallback == "" {
		fmt.Fprintf(w, "export KUBECOLOR_%s=%q\n",
			p.formatName(path),
			p.viper.GetString(strings.Join(path, ".")))
	} else {
		fmt.Fprintf(w, "export KUBECOLOR_%s=%q\n",
			p.formatName(path),
			fallback)
	}
	fmt.Fprintln(w)

	return nil
}

func (p *EnvPrinter) formatFallback(field Field) string {
	if field.DefaultFrom != "" {
		return fmt.Sprintf("$KUBECOLOR_%s", strings.ToUpper(strings.ReplaceAll(field.DefaultFrom, ".", "_")))
	}
	if field.DefaultFromMany != "" {
		split := strings.Split(field.DefaultFromMany, ",")
		for i, v := range split {
			split[i] = fmt.Sprintf("$KUBECOLOR_%s", strings.ToUpper(strings.ReplaceAll(v, ".", "_")))
		}
		return fmt.Sprintf("%s", strings.Join(split, "/"))
	}
	return ""
}

func (p *EnvPrinter) formatName(path []string) string {
	pathUppercase := slices.Clone(path)
	for i := range pathUppercase {
		pathUppercase[i] = strings.ToUpper(pathUppercase[i])
	}

	return strings.Join(pathUppercase, "_")
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
