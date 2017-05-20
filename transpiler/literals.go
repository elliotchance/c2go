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
	return &goast.BasicLit{
		Kind:  token.FLOAT,
		Value: fmt.Sprintf("%f", n.Value),
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

func transpilePredefinedExpr(n *ast.PredefinedExpr, p *program.Program) (*goast.BasicLit, string, error) {
	// A predefined expression is a literal that is not given a value until
	// compile time.
	//
	// TODO: Predefined expressions are not evaluated
	// https://github.com/elliotchance/c2go/issues/81

	var value string

	switch n.Name {
	case "__PRETTY_FUNCTION__":
		value = "[]byte(\"void print_number(int *)\")"

	case "__func__":
		value = fmt.Sprintf("[]byte(\"%s\")", p.FunctionName)

	default:
		// There are many more.
		panic(fmt.Sprintf("unknown PredefinedExpr: %s", n.Name))
	}

	return &goast.BasicLit{
		Kind:  token.STRING,
		Value: value,
	}, "const char*", nil
}
