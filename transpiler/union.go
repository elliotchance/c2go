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
)

func transpileUnion(name string, size int, fields []*goast.Field) (
	_ []goast.Decl, err error) {

	type field struct {
		Name      string
		TypeField string
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
	memory unsafe.Pointer
}

func (unionVar * {{ .Name }}) copy() ( {{ .Name }}){
	var buffer [{{ .Size }}]byte
	for i := range buffer{
		buffer[i] = (*((*[{{ .Size }}]byte)(unionVar.memory)))[i]
	}
	var newUnion {{ .Name }}
	newUnion.memory = unsafe.Pointer(&buffer)
	return newUnion
}

{{ range .Fields }}
func (unionVar * {{ $.Name }}) {{ .Name }}() (*{{ .TypeField }}){
	if unionVar.memory == nil{
		var buffer [{{ $.Size }}]byte
		unionVar.memory = unsafe.Pointer(&buffer)
	}
	return (*{{ .TypeField }})(unionVar.memory)
}
{{ end }}
`
	// Generate structure of union
	var un union
	un.Name = name
	un.Size = size
	for i := range fields {
		var f field
		f.Name = fields[i].Names[0].Name

		var buf bytes.Buffer
		err = format.Node(&buf, token.NewFileSet(), fields[i].Type)
		if err != nil {
			err = fmt.Errorf("cannot parse type '%s' : %v", fields[i].Type, err)
			return
		}
		f.TypeField = buf.String()

		un.Fields = append(un.Fields, f)
	}

	tmpl := template.Must(template.New("").Parse(src))
	var source bytes.Buffer
	err = tmpl.Execute(&source, un)
	if err != nil {
		err = fmt.Errorf("cannot execute template \"%s\" for data %v : %v", source.String(), un, err)
		return
	}

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", source.String(), 0)
	if err != nil {
		err = fmt.Errorf("cannot parse source \"%s\" : %v", source.String(), err)
		return
	}

	return f.Decls[1:], nil
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
