// This file contains utility and helper methods for making it easier to
// generate parts of the Go AST.

package util

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"regexp"
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

// IsAValidFunctionName performs a crude check to see if a string would make a
// valid function name in Go. Crude because it is not based on the actual Go
// standard and allows for some special conditions that allow writing code
// easier.
func IsAValidFunctionName(s string) bool {
	// This is a special case that is used by transpileBinaryExpression().
	if s == "(*[1]int)" {
		return true
	}

	// We allow '.' in the identifier name because there are lots of places
	// where we call a function like "os.Exit", but the dot is not strictly
	// allowed.
	//
	// The identifier may start with zero or more "[]" since types are used as
	// function names.
	return regexp.MustCompile(`^(\[\d*\])*[a-zA-Z_][a-zA-Z0-9_.]*$`).
		Match([]byte(s))
}

// NewCallExpr creates a new *"go/ast".CallExpr with each of the arguments
// (after the function name) being each of the expressions that represent the
// individual arguments.
//
// The function name is checked with IsAValidFunctionName and will panic if the
// function name is deemed to be not valid.
func NewCallExpr(functionName string, args ...goast.Expr) *goast.CallExpr {
	// Make sure the function name is valid. This can lead to some really
	// painful to debug errors. See Program.String().
	if !IsAValidFunctionName(functionName) {
		panic(fmt.Sprintf("not a valid function name: %s", functionName))
	}

	return &goast.CallExpr{
		Fun:  goast.NewIdent(functionName),
		Args: args,
	}
}

// NewFuncClosure creates a new *"go/ast".CallExpr that calls a function
// literal closure. The first argument is the Go return type of the
// closure, and the remainder of the arguments are the statements of the
// closure body.
func NewFuncClosure(returnType string, stmts ...goast.Stmt) *goast.CallExpr {
	return &goast.CallExpr{
		Fun: &goast.FuncLit{
			Type: &goast.FuncType{
				Params: &goast.FieldList{},
				Results: &goast.FieldList{
					List: []*goast.Field{
						&goast.Field{
							Type: goast.NewIdent(returnType),
						},
					},
				},
			},
			Body: &goast.BlockStmt{
				List: stmts,
			},
		},
		Args: []goast.Expr{},
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

func NewNil() *goast.BasicLit {
	return NewStringLit("nil")
}
