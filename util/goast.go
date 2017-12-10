// This file contains utility and helper methods for making it easier to
// generate parts of the Go AST.

package util

import (
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

// NewExprStmt returns a new ExprStmt from an expression. It is used when
// converting a single expression into a statement for another receiver.
//
// It is recommended you use this method of instantiating the ExprStmt yourself
// because NewExprStmt will check that the expr is not nil (or panic). This is
// much more helpful when trying to debug why the Go source build crashes
// because of a nil pointer - which eventually leads back to a nil expr.
func NewExprStmt(expr goast.Expr) *goast.ExprStmt {
	PanicIfNil(expr, "expr is nil")

	return &goast.ExprStmt{
		X: expr,
	}
}

// IsAValidFunctionName performs a check to see if a string would make a
// valid function name in Go. Go allows unicode characters, but C doesn't.
func IsAValidFunctionName(s string) bool {
	return GetRegex(`^[a-zA-Z_][a-zA-Z0-9_]*$`).
		Match([]byte(s))
}

// Convert a type as a string into a Go AST expression.
func typeToExpr(t string) goast.Expr {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("bad type: '%v'", t))
		}
	}()

	return internalTypeToExpr(t)
}

func internalTypeToExpr(t string) goast.Expr {
	// I'm not sure if this is an error or not. It is caused by processing the
	// resolved type of "void" which is "". It is used on functions to denote
	// that it does not have a return type.
	if t == "" {
		return nil
	}

	// Empty Interface
	if t == "interface{}" {
		return &goast.InterfaceType{
			Methods: &goast.FieldList{},
		}
	}

	// Parenthesis Expression
	if strings.HasPrefix(t, "(") && strings.HasSuffix(t, ")") {
		return &goast.ParenExpr{
			X: typeToExpr(t[1 : len(t)-1]),
		}
	}

	// Pointer Type
	if strings.HasPrefix(t, "*") {
		return &goast.StarExpr{
			X: typeToExpr(t[1:]),
		}
	}

	// Slice
	if strings.HasPrefix(t, "[]") {
		return &goast.ArrayType{
			Elt: typeToExpr(t[2:]),
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

	// Array access
	match := GetRegex(`(.+)\[(\d+)\]$`).FindStringSubmatch(t)
	if match != nil {
		return &goast.IndexExpr{
			X: typeToExpr(match[1]),
			// This should use NewIntLit, but it doesn't seem right to
			// cast the string to an integer to have it converted back to
			// as string.
			Index: &goast.BasicLit{
				Kind:  token.INT,
				Value: match[2],
			},
		}
	}

	// Fixed Length Array
	match = GetRegex(`^\[(\d+)\](.+)$`).FindStringSubmatch(t)
	if match != nil {
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

	// This may panic, and so it will be handled by typeToExpr().
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
			Type: NewFuncType(&goast.FieldList{}, returnType),
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
//
// Assignment operators in C can be nested inside other expressions, like:
//
//     a + (b += 3)
//
// In Go this is not allowed. Since the operators mutate variables it is not
// possible in some cases to move the statements before or after. The only safe
// (and generic) way around this is to create an immediately executing closure,
// like:
//
//     a + (func () int { b += 3; return b }())
//
// In a lot of cases this may be unnecessary and obfuscate the Go output but
// these will have to be optimised over time and be strict about the
// situation they are simplifying.
//
// If stmt is true then the binary expression is the whole statement. This means
// that the closure above does not need to applied. This makes the output code
// much neater.
func NewBinaryExpr(left goast.Expr, operator token.Token, right goast.Expr,
	returnType string, stmt bool) goast.Expr {
	PanicIfNil(left, "left is nil")
	PanicIfNil(right, "right is nil")

	var b goast.Expr = &goast.BinaryExpr{
		X:  left,
		Op: operator,
		Y:  right,
	}

	// Wrap the assignment operator in a closure if we must.
	if !stmt {
		switch operator {
		case token.ASSIGN,
			token.ADD_ASSIGN,
			token.SUB_ASSIGN,
			token.MUL_ASSIGN,
			token.QUO_ASSIGN,
			token.REM_ASSIGN,
			token.AND_ASSIGN,
			token.OR_ASSIGN,
			token.XOR_ASSIGN,
			token.SHL_ASSIGN,
			token.SHR_ASSIGN,
			token.AND_NOT_ASSIGN:
			returnStmt := &goast.ReturnStmt{
				Results: []goast.Expr{left},
			}

			b = NewFuncClosure(returnType, NewExprStmt(b), returnStmt)
		}
	}

	return b
}

// NewIdent - create a new Go ast Ident
func NewIdent(name string) *goast.Ident {
	// TODO: The name of a variable or field cannot be a reserved word
	// https://github.com/elliotchance/c2go/issues/83
	// Search for this issue in other areas of the codebase.
	if IsGoKeyword(name) {
		name += "_"
	}

	// Remove const prefix as it has no equivalent in Go.
	if strings.HasPrefix(name, "const ") {
		name = name[6:]
	}

	if !IsAValidFunctionName(name) {
		// Normally we do not panic because we want the transpiler to recover as
		// much as possible so that we always get Go output - even if it's
		// wrong. However, in this case we must panic because we know that this
		// identity will cause the AST renderer in Go to panic with a very
		// unhelpful error message.
		//
		// Panic now so that we can see where the bad identifier is coming from.
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

// NewIntLit - create a Go ast BasicLit for `INT` value
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
// instead of instantiating the UnaryExpr directly because this function has
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

// IsGoKeyword will return true if a word is one of the reserved words in Go.
// This means that it cannot be used as an identifier, function name, etc.
//
// The list of reserved words has been taken from the spec at
// https://golang.org/ref/spec#Keywords
func IsGoKeyword(w string) bool {
	switch w {
	case "break", "default", "func", "interface", "select", "case", "defer",
		"go", "map", "struct", "chan", "else", "goto", "package", "switch",
		"const", "fallthrough", "if", "range", "type", "continue", "for",
		"import", "return", "var":
		return true
	}

	return false
}

// ConvertFunctionNameFromCtoGo - convert function name fromC to Go
func ConvertFunctionNameFromCtoGo(name string) string {
	if name == "_" {
		return "__"
	}
	return name
}

func CreateSliceFromReference(goType string, expr goast.Expr) *goast.SliceExpr {
	// If the Go type is blank it means that the C type is 'void'.
	if goType == "" {
		goType = "interface{}"
	}

	// This is a hack to convert a reference to a variable into a slice that
	// points to the same location. It will look similar to:
	//
	//     (*[1]int)(unsafe.Pointer(&a))[:]
	//
	// You must always call this Go before using CreateSliceFromReference:
	//
	//     p.AddImport("unsafe")
	//
	return &goast.SliceExpr{
		X: NewCallExpr(
			fmt.Sprintf("(*[1]%s)", goType),
			NewCallExpr("unsafe.Pointer", &goast.UnaryExpr{
				X:  expr,
				Op: token.AND,
			}),
		),
	}
}

func NewFuncType(fieldList *goast.FieldList, returnType string) *goast.FuncType {
	returnTypes := []*goast.Field{}
	if returnType != "" {
		returnTypes = append(returnTypes, &goast.Field{
			Type: NewTypeIdent(returnType),
		})
	}

	return &goast.FuncType{
		Params: fieldList,
		Results: &goast.FieldList{
			List: returnTypes,
		},
	}
}

// NewGoExpr is used to facilitate the creation of go AST
//
// It should not be used to transpile user code.
func NewGoExpr(expr string) goast.Expr {
	e, err := parser.ParseExpr(expr)
	if err != nil {
		panic("programming error: " + expr)
	}

	return e
}
