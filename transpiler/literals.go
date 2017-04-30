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

func transpilePredefinedExpr(n *ast.PredefinedExpr, p *program.Program) (*goast.BasicLit, string, error) {
	// A predefined expression is a literal that is not given a value until
	// compile time.

	var value string

	switch n.Name {
	case "__PRETTY_FUNCTION__":
		// FIXME
		value = "\"void print_number(int *)\""

	case "__func__":
		// FIXME
		value = fmt.Sprintf("\"%s\"", "print_number")

	default:
		panic(fmt.Sprintf("unknown PredefinedExpr: %s", n.Name))
	}

	return &goast.BasicLit{
		Kind:  token.STRING,
		Value: value,
	}, "const char*", nil
}
