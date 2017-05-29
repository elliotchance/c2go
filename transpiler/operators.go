// This file contains functions transpiling some general operator expressions.
// See binary.go and unary.go.

package transpiler

import (
	"fmt"
	"go/token"

	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
)

// transpileConditionalOperator transpiles a conditional (also known as a
// ternary) operator:
//
//     a ? b : c
//
// We cannot simply convert these to an "if" statement becuase they by inside
// another expression.
//
// Since Go does not support the ternary operator or inline "if" statements we
// use a function, noarch.Ternary() to work the same way.
//
// It is also important to note that C only evaulates the "b" or "c" condition
// based on the result of "a" (from the above example). So we wrap the "b" and
// "c" in closures so that the Ternary function will only evaluate one of them.
func transpileConditionalOperator(n *ast.ConditionalOperator, p *program.Program) (
	*goast.CallExpr, string, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	a, _, newPre, newPost, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	b, _, newPre, newPost, err := transpileToExpr(n.Children[1], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	c, _, newPre, newPost, err := transpileToExpr(n.Children[2], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	p.AddImport("github.com/elliotchance/c2go/noarch")

	// The following code will generate the Go AST that will simulate a
	// conditional (ternary) operator, in the form of:
	//
	//     noarch.Ternary(
	//         $1,
	//         func () interface{} {
	//             return $2
	//         },
	//         func () interface{} {
	//             return $3
	//         },
	//     )
	//
	// $2 and $3 (the true and false condition respectively) must be wrapped in
	// a closure so that they are not both executed.
	return util.NewCallExpr(
		"noarch.Ternary",
		a,
		newTernaryWrapper(b),
		newTernaryWrapper(c),
	), n.Type, preStmts, postStmts, nil
}

// newTernaryWrapper is a helper method used by transpileConditionalOperator().
// It will wrap an expression in a closure.
func newTernaryWrapper(e goast.Expr) *goast.FuncLit {
	return &goast.FuncLit{
		Type: &goast.FuncType{
			Params: &goast.FieldList{},
			Results: &goast.FieldList{
				List: []*goast.Field{
					&goast.Field{
						Type: &goast.InterfaceType{
							Methods: &goast.FieldList{},
						},
					},
				},
			},
		},
		Body: &goast.BlockStmt{
			List: []goast.Stmt{
				&goast.ReturnStmt{
					Results: []goast.Expr{e},
				},
			},
		},
	}
}

// transpileParenExpr transpiles an expression that is wrapped in parentheses.
// There is a special case where "(0)" is treated as a NULL (since that's what
// the macro expands to). We have to return the type as "null" since we don't
// know at this point what the NULL expression will be used in conjuction with.
func transpileParenExpr(n *ast.ParenExpr, p *program.Program) (
	*goast.ParenExpr, string, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	e, eType, newPre, newPost, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	r := &goast.ParenExpr{
		X: e,
	}
	if types.IsNullExpr(r) {
		eType = "null"
	}

	return r, eType, preStmts, postStmts, nil
}

func transpileCompoundAssignOperator(n *ast.CompoundAssignOperator, p *program.Program) (
	*goast.BinaryExpr, string, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	left, _, newPre, newPost, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	right, _, newPre, newPost, err := transpileToExpr(n.Children[1], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	return &goast.BinaryExpr{
		X:  left,
		Y:  right,
		Op: getTokenForOperator(n.Opcode),
	}, "", preStmts, postStmts, nil
}

// getTokenForOperator returns the Go operator token for the provided C
// operator.
func getTokenForOperator(operator string) token.Token {
	switch operator {
	// Arithmetic
	case "--":
		return token.DEC
	case "++":
		return token.INC
	case "+":
		return token.ADD
	case "-":
		return token.SUB
	case "*":
		return token.MUL
	case "/":
		return token.QUO
	case "%":
		return token.REM

	// Assignment
	case "=":
		return token.ASSIGN
	case "+=":
		return token.ADD_ASSIGN
	case "-=":
		return token.SUB_ASSIGN
	case "*=":
		return token.MUL_ASSIGN
	case "/=":
		return token.QUO_ASSIGN
	case "%=":
		return token.REM_ASSIGN
	case "&=":
		return token.AND_ASSIGN
	case "|=":
		return token.OR_ASSIGN
	case "^=":
		return token.XOR_ASSIGN
	case "<<=":
		return token.SHL_ASSIGN
	case ">>=":
		return token.SHR_ASSIGN

	// Bitwise
	case "&":
		return token.AND
	case "|":
		return token.OR
	case "~":
		return token.XOR
	case ">>":
		return token.SHR
	case "<<":
		return token.SHL
	case "^":
		return token.XOR

	// Comparison
	case ">=":
		return token.GEQ
	case "<=":
		return token.LEQ
	case "<":
		return token.LSS
	case ">":
		return token.GTR
	case "!=":
		return token.NEQ
	case "==":
		return token.EQL

	// Logical
	case "!":
		return token.NOT
	case "&&":
		return token.LAND
	case "||":
		return token.LOR

	// Other
	case ",":
		return token.COMMA
	}

	panic(fmt.Sprintf("unknown operator: %s", operator))
}
