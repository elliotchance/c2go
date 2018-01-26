package program

import (
	"fmt"
	"reflect"

	"github.com/elliotchance/c2go/ast"
)

// GenerateErrorMessage - generate error message
func (p *Program) GenerateErrorMessage(e error, n ast.Node) string {
	if e != nil {
		structName := reflect.TypeOf(n).Elem().Name()
		return fmt.Sprintf("// Error (%s): %s: %s", structName,
			n.Position().GetSimpleLocation(), e.Error())
	}

	return ""
}

// GenerateWarningMessage - generate warning message
func (p *Program) GenerateWarningMessage(e error, n ast.Node) string {
	if e != nil {
		structName := reflect.TypeOf(n).Elem().Name()
		return fmt.Sprintf("// Warning (%s): %s: %s", structName,
			n.Position().GetSimpleLocation(), e.Error())
	}

	return ""
}

// GenerateWarningOrErrorMessage - generate error if it happen
func (p *Program) GenerateWarningOrErrorMessage(e error, n ast.Node, isError bool) string {
	if isError {
		return p.GenerateErrorMessage(e, n)
	}

	return p.GenerateWarningMessage(e, n)
}
