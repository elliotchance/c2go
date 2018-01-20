package transpiler

import (
	"bytes"
	"html/template"

	goast "go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/elliotchance/c2go/util"
)

func transpileUnion(name string, size int, fields []*goast.Field) (
	_ []goast.Decl, err error) {

	type field struct {
		Name          string
		PositionField int
		TypeField     string
	}

	type union struct {
		Name   string
		Size   int
		Fields []field
	}

	src := `package main

import(
	"unsafe"
	"reflect"
)

type {{ .Name }} struct{
	memory [{{ .Size }}]byte
	pointer interface{}
}

`
	// Generate structure of union
	var un union
	un.Name = name
	un.Size = size
	for i := range fields {
		var f field
		f.Name = fields[i].Names[0].Name

		var buf bytes.Buffer
		err := format.Node(&buf, token.NewFileSet(), fields[i].Type)
		if err != nil {
			panic(err)
		}
		f.TypeField = buf.String()

		f.PositionField = i

		un.Fields = append(un.Fields, f)
	}

	tmpl := template.Must(template.New("").Parse(src))
	var source bytes.Buffer
	err = tmpl.Execute(&source, un)
	if err != nil {
		return
	}

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", source.String(), 0)
	if err != nil {
		return
	}

	return f.Decls[1:], nil
}

func getUnionVariable(goType string, union goast.Expr) goast.Expr {
	return &goast.StarExpr{
		X: &goast.CallExpr{
			Fun: &goast.ParenExpr{
				Lparen: 1,
				X: &goast.StarExpr{
					X: goast.NewIdent(goType),
					//X: &goast.ArrayType{Elt: goast.NewIdent(goType)},
				},
			},
			Lparen: 1,
			Args: []goast.Expr{&goast.CallExpr{
				Fun: &goast.SelectorExpr{
					X:   goast.NewIdent("unsafe"),
					Sel: goast.NewIdent("Pointer"),
				},
				Lparen: 1,
				Args: []goast.Expr{
					&goast.UnaryExpr{
						Op: token.AND,
						X: &goast.SelectorExpr{
							X:   union,
							Sel: util.NewIdent("memory"),
						},
					},
				},
			}},
		},
	}

}
