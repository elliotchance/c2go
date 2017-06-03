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

// NewExprStmt returns a new ExprStmt from an expression. It is used when
// converting a single expression into a statement for another receiver.
//
// It is recommended you use this method of instantiating the ExprStmt yourself
// because NewExprStmt will check that the expr is not nil (or panic). This is
// much more helpful when trying to debug why the Go source build crashes
// becuase of a nil pointer - which eventually leads back to a nil expr.
func NewExprStmt(expr goast.Expr) *goast.ExprStmt {
	PanicIfNil(expr, "expr is nil")

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

// NewBinaryExpr create a new Go AST binary expression with a left, operator and
// right operand.
//
// You should use this instead of BinaryExpr directly so that nil left and right
// operands can be caught (and panic) before Go tried to render the source -
// which would result in a very hard to debug error.
func NewBinaryExpr(left goast.Expr, operator token.Token, right goast.Expr) *goast.BinaryExpr {
	PanicIfNil(left, "left is nil")
	PanicIfNil(right, "right is nil")

	return &goast.BinaryExpr{
		X:  left,
		Op: operator,
		Y:  right,
	}
}

// NewStringLit returns a new Go basic literal with a string value.
func NewStringLit(value string) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.STRING,
		Value: value,
	}
}

// NewIntLit returns a new Go basic literal with an integer value.
func NewIntLit(value int) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.INT,
		Value: strconv.Itoa(value),
	}
}

// NewNil returns a Go AST identity that can be used to represent "nil".
func NewNil() *goast.Ident {
	return goast.NewIdent("nil")
}
