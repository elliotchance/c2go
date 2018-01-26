package transpiler

import (
	"bytes"
	"html/template"
	"strings"

	goast "go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"fmt"
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
	value    interface{}
	arr      [{{ .Size }}]byte
}

func (self *{{ .Name }}) cast(t reflect.Type) reflect.Value {
	return reflect.NewAt(t, unsafe.Pointer(&self.arr[0]))
}

func (self *{{ .Name }}) assign(v interface{}){
	value := reflect.ValueOf(v).Elem()
	value.Set(self.cast(value.Type()).Elem())
}

func (self *{{ .Name }}) UntypedSet(v interface{}){
	value := reflect.ValueOf(v)
	self.cast(value.Type()).Elem().Set(value)
}

{{ range .Fields }}
// Get{{ .Name }} - return value of {{ .Name }}
func (self *{{ $.Name }}) Get{{ .Name }} () (res {{ .TypeField }}){
	self.assign(&res)
	return
}

// Set{{ .Name }} - set value of {{ .Name }}
func (self *{{ $.Name }}) Set{{ .Name }} (v {{ .TypeField }}) {{ .TypeField }}{
	self.value = v // added for avoid GC removing pointers in union
	self.UntypedSet(v)
	return v
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
			panic(err)
		}
		f.TypeField = buf.String()

		// capitalization first letter
		name := strings.ToUpper(string(f.Name[0]))
		if len(f.Name) > 1 {
			name += f.Name[1:]
		}
		f.Name = name

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

func getFunctionNameForUnion(verb, variableName, variableType, attributeName string) string {
	if strings.HasPrefix(variableType, "[]") {
		return fmt.Sprintf("%s[0].%s%s", variableName, verb, strings.Title(attributeName))
	}

	return fmt.Sprintf("%s.%s%s", variableName, verb, strings.Title(attributeName))
}

func getFunctionNameForUnionGetter(variableName, variableType, attributeName string) string {
	return getFunctionNameForUnion("Get", variableName, variableType, attributeName)
}

func getFunctionNameForUnionSetter(variableName, variableType, attributeName string) string {
	return getFunctionNameForUnion("Set", variableName, variableType, attributeName)
}
