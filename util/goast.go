// This file contains utility and helper methods for making it easier to
// generate parts of the Go AST.

package util

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"
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

// IsAValidFunctionName performs a check to see if a string would make a
// valid function name in Go. Go allows unicode characters, but C doesn't.
func IsAValidFunctionName(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).
		Match([]byte(s))
}

// Convert a type as a string into a Go AST expression
func typeToExpr(t string) goast.Expr {

	// Empty Interface
	if t == "interface{}" {
		return &goast.InterfaceType{Methods: &goast.FieldList{}}
	}

	// Parenthesis Expression
	if strings.HasPrefix(t, "(") && strings.HasSuffix(t, ")") {
		return &goast.ParenExpr{X: typeToExpr(t[1 : len(t)-1])}
	}

	// Pointer Type
	if strings.HasPrefix(t, "*") {
		return &goast.StarExpr{X: typeToExpr(t[1:])}
	}

	// Slice
	if strings.HasPrefix(t, "[]") {
		return &goast.ArrayType{Elt: typeToExpr(t[2:])}
	}

	// Fixed Length Array
	if match := regexp.MustCompile(`^\[(\d+)\](.+)$`).FindStringSubmatch(t); match != nil {
		return &goast.ArrayType{
			Elt: typeToExpr(match[2]),
			// This should use NewIntLit, but it doesn't seem right to
			// cast the string to an integer to have it converted back to
			// as string.
			Len: &goast.BasicLit{
				Kind:  token.INT,
				Value: match[1],
			},
		}
	}

	// Selector: "type.identifier"
	if strings.Contains(t, ".") {
		i := strings.IndexByte(t, '.')
		return &goast.SelectorExpr{
			X:   typeToExpr(t[0:i]),
			Sel: NewIdent(t[i+1:]),
		}
	}

	return NewIdent(t)

}

// NewCallExpr creates a new *"go/ast".CallExpr with each of the arguments
// (after the function name) being each of the expressions that represent the
// individual arguments.
//
// The function name is checked with IsAValidFunctionName and will panic if the
// function name is deemed to be not valid.
func NewCallExpr(functionName string, args ...goast.Expr) *goast.CallExpr {
	return &goast.CallExpr{
		Fun:  typeToExpr(functionName),
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
							Type: NewTypeIdent(returnType),
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

func NewIdent(name string) *goast.Ident {
	if !IsAValidFunctionName(name) {
		panic(fmt.Sprintf("invalid identity: '%s'", name))
	}

	return goast.NewIdent(name)
}

// NewTypeIdent created a new Go identity that is to be used for a Go type. This
// is different from NewIdent in how the input string is validated.
func NewTypeIdent(name string) goast.Expr {
	return typeToExpr(name)
}

// NewStringLit returns a new Go basic literal with a string value.
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

// NewFloatLit creates a new Float Literal.
func NewFloatLit(value float64) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.FLOAT,
		Value: strconv.FormatFloat(value, 'g', -1, 64),
	}
}

// NewNil returns a Go AST identity that can be used to represent "nil".
func NewNil() *goast.Ident {
	return NewIdent("nil")
}

// NewUnaryExpr creates a new Go unary expression. You should use this function
// instead of instantiating the UnaryExpr directly because this funtion has
// extra error checking.
func NewUnaryExpr(operator token.Token, right goast.Expr) *goast.UnaryExpr {
	if right == nil {
		panic("right is nil")
	}

	return &goast.UnaryExpr{
		Op: operator,
		X:  right,
	}
}
