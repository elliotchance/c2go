package program

import (
	"fmt"
	"github.com/elliotchance/c2go/ast"
	"reflect"
	"regexp"
)

// getNicerLineNumber tries to extract a more useful line number from a
// position. If the line number cannot be determined then the original location
// string is returned.
func getNicerLineNumber(s string) string {
	matches := regexp.MustCompile(`line:(\d+)`).FindStringSubmatch(s)
	if len(matches) > 0 {
		return fmt.Sprintf("line %s", matches[1])
	}

	return s
}

func (p *Program) GenerateErrorMessage(e error, n ast.Node) string {
	if e != nil {
		structName := reflect.TypeOf(n).Elem().Name()
		return fmt.Sprintf("// Error (%s): %s: %s", structName,
			getNicerLineNumber(ast.Position(n)), e.Error())
	}

	return ""
}

func (p *Program) GenerateWarningMessage(e error, n ast.Node) string {
	if e != nil {
		structName := reflect.TypeOf(n).Elem().Name()
		return fmt.Sprintf("// Warning (%s): %s: %s", structName,
			getNicerLineNumber(ast.Position(n)), e.Error())
	}

	return ""
}

func (p *Program) GenerateWarningOrErrorMessage(e error, n ast.Node, isError bool) string {
	if isError {
		return p.GenerateErrorMessage(e, n)
	}

	return p.GenerateWarningMessage(e, n)
}
