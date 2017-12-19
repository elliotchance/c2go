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
	_ *goast.CallExpr, theType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile ConditionalOperator : err = %v", err)
		}
	}()

	// a - condition
	a, aType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// null in C is zero
	if aType == types.NullPointer {
		a = &goast.BasicLit{
			Kind:  token.INT,
			Value: "0",
		}
		aType = "int"
	}

	a, err = types.CastExpr(p, a, aType, "bool")
	if err != nil {
		return
	}

	// b - body
	b, bType, newPre, newPost, err := transpileToExpr(n.Children()[1], p, false)
	if err != nil {
		return
	}
	// Theorephly, lenght is must be zero
	if len(newPre) > 0 || len(newPost) > 0 {
		p.AddMessage(p.GenerateWarningMessage(
			fmt.Errorf("lenght of pre or post in body must be zero. {%d,%d}", len(newPre), len(newPost)), n))
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	if n.Type != "void" {
		b, err = types.CastExpr(p, b, bType, n.Type)
		if err != nil {
			return
		}
		bType = n.Type
	}

	// c - else body
	c, cType, newPre, newPost, err := transpileToExpr(n.Children()[2], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	if n.Type != "void" {
		c, err = types.CastExpr(p, c, cType, n.Type)
		if err != nil {
			return
		}
		cType = n.Type
	}

	// rightType - generate return type
	var returnType string
	if n.Type != "void" {
		returnType, err = types.ResolveType(p, n.Type)
		if err != nil {
			return
		}
	}

	var bod, els goast.BlockStmt

	bod.Lbrace = 1
	if bType != types.ToVoid {
		if n.Type != "void" {
			bod.List = []goast.Stmt{
				&goast.ReturnStmt{
					Results: []goast.Expr{b},
				},
			}
		} else {
			bod.List = []goast.Stmt{&goast.ExprStmt{b}}
		}
	}

	els.Lbrace = 1
	if cType != types.ToVoid {
		if n.Type != "void" {
			els.List = []goast.Stmt{
				&goast.ReturnStmt{
					Results: []goast.Expr{c},
				},
			}
		} else {
			els.List = []goast.Stmt{&goast.ExprStmt{c}}
		}
	}

	return util.NewFuncClosure(
		returnType,
		&goast.IfStmt{
			Cond: a,
			Body: &bod,
			Else: &els,
		},
	), n.Type, preStmts, postStmts, nil
}

// transpileParenExpr transpiles an expression that is wrapped in parentheses.
// There is a special case where "(0)" is treated as a NULL (since that's what
// the macro expands to). We have to return the type as "null" since we don't
// know at this point what the NULL expression will be used in conjunction with.
func transpileParenExpr(n *ast.ParenExpr, p *program.Program) (
	r *goast.ParenExpr, exprType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile ParenExpr. err = %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	expr, exprType, preStmts, postStmts, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return
	}

	if exprType == types.NullPointer {
		r = &goast.ParenExpr{X: expr}
		return
	}

	if !types.IsFunction(exprType) && exprType != "void" && exprType != types.ToVoid {
		expr, err = types.CastExpr(p, expr, exprType, n.Type)
		if err != nil {
			return
		}
		exprType = n.Type
	}

	r = &goast.ParenExpr{X: expr}

	return
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
