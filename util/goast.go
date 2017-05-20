// This file contains utility and helper methods for making it easier to
// generate parts of the Go AST.

package util

import (
	goast "go/ast"
	"go/token"
	"strconv"
)

func NewExprStmt(expr goast.Expr) *goast.ExprStmt {
	if expr == nil {
		panic("expr is nil")
	}

	return &goast.ExprStmt{
		X: expr,
	}
}

func NewCallExpr(functionName string, args ...goast.Expr) *goast.CallExpr {
	return &goast.CallExpr{
		Fun:  goast.NewIdent(functionName),
		Args: args,
	}
}

func NewBinaryExpr(left goast.Expr, operator token.Token, right goast.Expr) *goast.BinaryExpr {
	if left == nil {
		panic("left is nil")
	}
	if right == nil {
		panic("right is nil")
	}

	return &goast.BinaryExpr{
		X:  left,
		Op: operator,
		Y:  right,
	}
}

func NewIdent(name string) goast.Expr {
	return goast.NewIdent(name)
}

func NewIdents(names ...string) []goast.Expr {
	idents := []goast.Expr{}

	for _, name := range names {
		idents = append(idents, goast.NewIdent(name))
	}

	return idents
}

func NewStringLit(value string) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.STRING,
		Value: value,
	}
}

func NewIntLit(value int) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.INT,
		Value: strconv.Itoa(value),
	}
}
