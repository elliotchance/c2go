package program

import (
	"fmt"

	"github.com/elliotchance/c2go/ast"
)

// Struct represents the definition for a C struct.
type Struct struct {
	// The name of the struct.
	Name string

	// Each of the fields and their C type. The field may be a string or an
	// instance of Struct for nested structures.
	Fields map[string]interface{}
}

// NewStruct creates a new Struct definition from an ast.RecordDecl.
func NewStruct(n *ast.RecordDecl) *Struct {
	fields := make(map[string]interface{})

	for _, field := range n.Children {
		switch f := field.(type) {
		case *ast.FieldDecl:
			fields[f.Name] = f.Type

		case *ast.RecordDecl:
			fields[f.Name] = NewStruct(f)

		case *ast.MaxFieldAlignmentAttr, *ast.AlignedAttr:
			// FIXME: Should these really be ignored?

		default:
			panic(fmt.Sprintf("cannot decode: %#v", f))
		}
	}

	return &Struct{
		Name:   n.Name,
		Fields: fields,
	}
}
