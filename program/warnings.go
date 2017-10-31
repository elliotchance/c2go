package program

import (
	"fmt"
	"reflect"

	"github.com/elliotchance/c2go/ast"
)

func (p *Program) GenerateErrorMessage(e error, n ast.Node) string {
	if e != nil {
		structName := reflect.TypeOf(n).Elem().Name()
		return fmt.Sprintf("// Error (%s): %d: %s", structName,
			n.Position().Line, e.Error())
	}

	return ""
}

func (p *Program) GenerateWarningMessage(e error, n ast.Node) string {
	if e != nil {
		structName := reflect.TypeOf(n).Elem().Name()
		return fmt.Sprintf("// Warning (%s): %d: %s", structName,
			n.Position().Line, e.Error())
	}

	return ""
}

func (p *Program) GenerateWarningOrErrorMessage(e error, n ast.Node, isError bool) string {
	if isError {
		return p.GenerateErrorMessage(e, n)
	}

	return p.GenerateWarningMessage(e, n)
}
