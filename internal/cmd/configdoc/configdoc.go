package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	md := MarkdownPrinter{Program: p}
	if err := md.Print(); err != nil {
		return fmt.Errorf("markdown: %w", err)
	}
	fmt.Println()
	yaml := YAMLPrinter{Program: p}
	if err := yaml.Print(); err != nil {
		return fmt.Errorf("yaml: %w", err)
	}
	return nil
}

func (p *Program) findCategory(typeName string) (Category, bool) {
	for _, c := range p.categories {
		if c.Type == typeName {
			return c, true
		}
	}
	return Category{}, false
}
