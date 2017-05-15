package program

import (
	"github.com/elliotchance/c2go/ast"
)

// Struct represents the definition for a C struct.
type Struct struct {
	// The name of the struct.
	Name string

	// Each of the fields and their C type.
	Fields map[string]string
}

// NewStruct creates a new Struct definition from an ast.RecordDecl.
func NewStruct(n *ast.RecordDecl) Struct {
	fields := make(map[string]string)

	for _, field := range n.Children {
		f := field.(*ast.FieldDecl)
		fields[f.Name] = f.Type
	}

	return Struct{
		Name:   n.Name,
		Fields: fields,
	}
}
