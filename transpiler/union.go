package transpiler

import (
	"bytes"
	"fmt"
	"html/template"

	goast "go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
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

// ((*(*[100]uint8)(unsafe.Pointer(&sha.u.memory)))[:])[0]
// ((*(*[ 25]int  )(unsafe.Pointer(&sha.u.memory)))[:])[0]
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

func isUnionMemberExpr(p *program.Program, n *ast.MemberExpr) (IsUnion bool) {
	if len(n.Children()) > 0 {
		if v, ok := n.Children()[0].(*ast.MemberExpr); ok {
			if p.IsUnion(v.Type) {
				IsUnion = true
			}
		}
		if v, ok := n.Children()[0].(*ast.DeclRefExpr); ok {
			if p.IsUnion(v.Type) {
				IsUnion = true
			}
		}
		if v, ok := n.Children()[0].(*ast.ImplicitCastExpr); ok {
			if p.IsUnion(v.Type) {
				IsUnion = true
			}
		}
	}
	return
}

func unionVariable(p *program.Program, n *ast.MemberExpr, x goast.Expr) (
	_ goast.Expr, cType string, ok bool) {
	if isUnionMemberExpr(p, n) {
		cType := n.Type
		var goType string
		var err error
		if types.IsFunction(cType) {
			goType, err = types.ResolveFunction(p, cType)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		} else {
			goType, err = types.ResolveType(p, cType)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
		return getUnionVariable(goType, x),
			n.Type, true
	}
	panic(fmt.Errorf("That MemberExpr is not union"))
}
