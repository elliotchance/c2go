// This file contains transpiling functions for literals and constants. Literals
// are single values like 123 or "hello".

package transpiler

import (
	"fmt"
	"go/token"

	goast "go/ast"

	"strconv"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
)

func transpileFloatingLiteral(n *ast.FloatingLiteral) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.FLOAT,
		Value: fmt.Sprintf("%f", n.Value),
	}
}

func transpileStringLiteral(n *ast.StringLiteral) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(n.Value),
	}
}

func transpileIntegerLiteral(n *ast.IntegerLiteral) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.INT,
		Value: strconv.Itoa(n.Value),
	}
}

func transpileCharacterLiteral(n *ast.CharacterLiteral) *goast.BasicLit {
	var s string

	// TODO: Transpile special character literals
	// https://github.com/elliotchance/c2go/issues/80
	switch n.Value {
	case '\n':
		s = "\\n"
	default:
		s = fmt.Sprintf("%c", n.Value)
	}

	return &goast.BasicLit{
		Kind:  token.CHAR,
		Value: fmt.Sprintf("'%s'", s),
	}
}

func transpilePredefinedExpr(n *ast.PredefinedExpr, p *program.Program) (*goast.BasicLit, string, error) {
	// A predefined expression is a literal that is not given a value until
	// compile time.
	//
	// TODO: Predefined expressions are not evaluated
	// https://github.com/elliotchance/c2go/issues/81

	var value string

	switch n.Name {
	case "__PRETTY_FUNCTION__":
		value = "\"void print_number(int *)\""

	case "__func__":
		value = fmt.Sprintf("\"%s\"", "print_number")

	default:
		// There are many more.
		panic(fmt.Sprintf("unknown PredefinedExpr: %s", n.Name))
	}

	return &goast.BasicLit{
		Kind:  token.STRING,
		Value: value,
	}, "const char*", nil
}
