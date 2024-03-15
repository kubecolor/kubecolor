package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
)

type Category struct {
	Type   string
	Fields []Field
}

type Field struct {
	Name    string
	Type    string
	Comment string

	DefaultFrom     string
	DefaultFromMany string
}

func ParseCategories() ([]Category, error) {
	b, err := os.ReadFile(flags.file)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "theme.go", b, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	categories := make([]Category, 0, len(f.Decls))
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		category, ok := visitGenDecl(genDecl)
		if !ok {
			continue
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func visitGenDecl(decl *ast.GenDecl) (Category, bool) {
	typeSpec := getTypeSpec(decl)
	if typeSpec == nil {
		return Category{}, false
	}
	if !strings.HasPrefix(typeSpec.Name.Name, "Theme") {
		return Category{}, false
	}
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return Category{}, false
	}
	category := Category{
		Type: typeSpec.Name.Name,
	}
	for _, field := range structType.Fields.List {
		if len(field.Names) != 1 {
			continue
		}
		name := field.Names[0].Name
		typ, ok := field.Type.(*ast.Ident)
		if !ok {
			continue
		}
		var comments []string
		if field.Comment != nil {
			for _, c := range field.Comment.List {
				comments = append(comments, trimComment(c.Text))
			}
		}

		var tag reflect.StructTag
		if field.Tag != nil {
			tag = reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
		}

		category.Fields = append(category.Fields, Field{
			Name:    name,
			Type:    typ.Name,
			Comment: strings.Join(comments, "\n"),

			DefaultFrom:     tag.Get("defaultFrom"),
			DefaultFromMany: tag.Get("defaultFromMany"),
		})
	}
	return category, true
}

func trimComment(s string) string {
	if after, ok := strings.CutPrefix(s, "//"); ok {
		return strings.TrimSpace(after)
	}
	if after, ok := strings.CutPrefix(s, "/*"); ok {
		inBetween, _ := strings.CutSuffix(after, "*/")
		return strings.TrimSpace(inBetween)
	}
	return s
}

func getTypeSpec(decl *ast.GenDecl) *ast.TypeSpec {
	if decl.Tok.String() != "type" {
		return nil
	}
	if len(decl.Specs) != 1 {
		return nil
	}
	spec, ok := decl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil
	}
	return spec
}
