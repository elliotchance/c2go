// This file contains transpiling functions for literals and constants. Literals
// are single values like 123 or "hello".

package transpiler

import (
	"bytes"
	"fmt"
	"go/token"

	goast "go/ast"

	"strconv"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
)

func transpileFloatingLiteral(n *ast.FloatingLiteral) *goast.BasicLit {
	return util.NewFloatLit(n.Value)
}

func transpileStringLiteral(n *ast.StringLiteral) goast.Expr {
	// Example:
	// StringLiteral 0x280b918 <col:29> 'char [30]' lvalue "%0"
	s, err := types.GetAmountArraySize(n.Type)
	if err != nil {
		return util.NewCallExpr("[]byte",
			util.NewStringLit(strconv.Quote(n.Value+"\x00")))
	}
	buf := bytes.NewBufferString(n.Value + "\x00")
	if buf.Len() < s {
		buf.Write(make([]byte, s-buf.Len()))
	}
	return util.NewCallExpr("[]byte",
		util.NewStringLit(strconv.Quote(buf.String())))
}

func transpileIntegerLiteral(n *ast.IntegerLiteral) (ret goast.Expr) {
	ret = &goast.BasicLit{
		Kind:  token.INT,
		Value: n.Value,
	}
	if n.Type == "int" {
		ret = util.NewCallExpr("int32", ret)
	}
	return
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

func transpileCompoundLiteralExpr(n *ast.CompoundLiteralExpr, p *program.Program) (goast.Expr, string, error) {
	expr, t, _, _, err := transpileToExpr(n.Children()[0], p, false)
	return expr, t, err
}
