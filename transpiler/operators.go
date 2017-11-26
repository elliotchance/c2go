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
// We cannot simply convert these to an "if" statement because they by inside
// another expression.
//
// Since Go does not support the ternary operator or inline "if" statements we
// use a closure to work the same way.
//
// It is also important to note that C only evaulates the "b" or "c" condition
// based on the result of "a" (from the above example).
func transpileConditionalOperator(n *ast.ConditionalOperator, p *program.Program) (
	*goast.CallExpr, string, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	a, aType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	b, bType, newPre, newPost, err := transpileToExpr(n.Children()[1], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	c, cType, newPre, newPost, err := transpileToExpr(n.Children()[2], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	a, err = types.CastExpr(p, a, aType, "bool")
	if err != nil {
		return nil, "", nil, nil, err
	}

	// TODO: Here it is being assumed that the return type of the
	// conditional operator is the type of the 'false' result. Things
	// are a bit more complicated then that in C.

	b, err = types.CastExpr(p, b, bType, cType)
	if err != nil {
		return nil, "", nil, nil, err
	}

	returnType, err := types.ResolveType(p, cType)
	if err != nil {
		return nil, "", nil, nil, err
	}

	return util.NewFuncClosure(
		returnType,
		&goast.IfStmt{
			Cond: a,
			Body: &goast.BlockStmt{
				List: []goast.Stmt{
					&goast.ReturnStmt{
						Results: []goast.Expr{b},
					},
				},
			},
			Else: &goast.BlockStmt{
				List: []goast.Stmt{
					&goast.ReturnStmt{
						Results: []goast.Expr{c},
					},
				},
			},
		},
	), cType, preStmts, postStmts, nil
}

// transpileParenExpr transpiles an expression that is wrapped in parentheses.
// There is a special case where "(0)" is treated as a NULL (since that's what
// the macro expands to). We have to return the type as "null" since we don't
// know at this point what the NULL expression will be used in conjunction with.
func transpileParenExpr(n *ast.ParenExpr, p *program.Program) (
	*goast.ParenExpr, string, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	e, eType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
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

func transpileCompoundAssignOperator(n *ast.CompoundAssignOperator, p *program.Program, exprIsStmt bool) (
	goast.Expr, string, []goast.Stmt, []goast.Stmt, error) {
	operator := getTokenForOperator(n.Opcode)
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	right, rightType, newPre, newPost, err := transpileToExpr(n.Children()[1], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// Construct code for computing compound assign operation to an union field
	memberExpr, ok := n.Children()[0].(*ast.MemberExpr)
	if ok {
		ref := memberExpr.GetDeclRefExpr()
		if ref != nil {
			// Get operator by removing last char that is '=' (e.g.: += becomes +)
			binaryOperation := n.Opcode
			binaryOperation = binaryOperation[:(len(binaryOperation) - 1)]

			// TODO: Is this duplicate code in unary.go?
			union := p.GetStruct(ref.Type)
			if union != nil && union.IsUnion {
				attrType, err := types.ResolveType(p, ref.Type)
				if err != nil {
					p.AddMessage(p.GenerateWarningMessage(err, memberExpr))
				}

				// Method names
				getterName := getFunctionNameForUnionGetter(ref.Name, attrType, memberExpr.Name)
				setterName := getFunctionNameForUnionSetter(ref.Name, attrType, memberExpr.Name)

				// Call-Expression argument
				argLHS := util.NewCallExpr(getterName)
				argOp := getTokenForOperator(binaryOperation)
				argRHS := right
				argValue := util.NewBinaryExpr(argLHS, argOp, argRHS, "interface{}", exprIsStmt)

				// Make Go expression
				resExpr := util.NewCallExpr(setterName, argValue)

				return resExpr, "", preStmts, postStmts, nil
			}
		}
	}

	left, leftType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// The right hand argument of the shift left or shift right operators
	// in Go must be unsigned integers. In C, shifting with a negative shift
	// count is undefined behaviour (so we should be able to ignore that case).
	// To handle this, cast the shift count to a uint64.
	if operator == token.SHL_ASSIGN || operator == token.SHR_ASSIGN {
		right, err = types.CastExpr(p, right, rightType, "unsigned long long")
		p.AddMessage(p.GenerateWarningOrErrorMessage(err, n, right == nil))
		if right == nil {
			right = util.NewNil()
		}
	}

	resolvedLeftType, err := types.ResolveType(p, leftType)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
	}

	return util.NewBinaryExpr(left, operator, right, resolvedLeftType, exprIsStmt),
		"", preStmts, postStmts, nil
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
