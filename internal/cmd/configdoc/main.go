package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/kubecolor/kubecolor/config"
	"github.com/spf13/viper"
)

var flags = struct {
	file    string
	printer string
}{
	file:    filepath.Join("config", "theme.go"),
	printer: "all",
}

func init() {
	flag.StringVar(&flags.file, "file", flags.file, "Path to theme.go file")
	flag.StringVar(&flags.printer, "printer", flags.printer, "What to print. One of: all, file, env")
}

func main() {
	flag.Parse()

	prog := Program{}
	if err := prog.Run(); err != nil {
		slog.Error("Error: " + err.Error())
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
	if flags.printer == "all" || flags.printer == "env" {
		md := EnvPrinter{Program: p}
		if err := md.Print(); err != nil {
			return fmt.Errorf("markdown: %w", err)
		}
	}
	if flags.printer == "all" {
		fmt.Println()
	}
	if flags.printer == "all" || flags.printer == "file" {
		yaml := FilePrinter{Program: p}
		if err := yaml.Print(); err != nil {
			return fmt.Errorf("yaml: %w", err)
		}
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
