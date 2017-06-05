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
	"github.com/elliotchance/c2go/util"
)

func transpileFloatingLiteral(n *ast.FloatingLiteral) *goast.BasicLit {
	// Important: Using 'e' as the formatter will output scientific notation.
	//
	// This is ugly in a lot of scenarios but it is more correct when dealing
	// with extremely large or small numbers that cannot be easily represented
	// as a decimal.
	//
	// This function should be improved to make smarter decisions about the best
	// formatting based on the size of the number.
	return &goast.BasicLit{
		Kind:  token.FLOAT,
		Value: fmt.Sprintf("%e", n.Value),
	}
}

func transpileStringLiteral(n *ast.StringLiteral) goast.Expr {
	return util.NewCallExpr("[]byte",
		util.NewStringLit(strconv.Quote(n.Value+"\x00")))
}

func transpileIntegerLiteral(n *ast.IntegerLiteral) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.INT,
		Value: n.Value,
	}
}

func transpileCharacterLiteral(n *ast.CharacterLiteral) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.CHAR,
		Value: fmt.Sprintf("%q", n.Value),
	}
}

func transpilePredefinedExpr(n *ast.PredefinedExpr, p *program.Program) (goast.Expr, string, error) {
	// A predefined expression is a literal that is not given a value until
	// compile time.
	//
	// TODO: Predefined expressions are not evaluated
	// https://github.com/elliotchance/c2go/issues/81

	switch n.Name {
	case "__PRETTY_FUNCTION__":
		return util.NewCallExpr(
			"[]byte",
			util.NewStringLit(`"void print_number(int *)"`),
		), "const char*", nil

	case "__func__":
		return util.NewCallExpr(
			"[]byte",
			util.NewStringLit(strconv.Quote(p.Function.Name)),
		), "const char*", nil

	default:
		// There are many more.
		panic(fmt.Sprintf("unknown PredefinedExpr: %s", n.Name))
	}
}
