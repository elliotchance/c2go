// This file contains functions for transpiling binary operator expressions.

package transpiler

import (
	goast "go/ast"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
)

func transpileBinaryOperator(n *ast.BinaryOperator, p *program.Program) (
	*goast.BinaryExpr, string, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}
	var err error

	left, leftType, newPre, newPost, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	right, rightType, newPre, newPost, err := transpileToExpr(n.Children[1], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	operator := getTokenForOperator(n.Operator)
	returnType := types.ResolveTypeForBinaryOperator(p, n.Operator, leftType, rightType)

	if operator == token.LAND {
		left, err = types.CastExpr(p, left, leftType, "bool")
		ast.WarningOrError(err, n, left == nil)
		if left == nil {
			left = util.NewStringLit("nil")
		}

		right, err = types.CastExpr(p, right, rightType, "bool")
		ast.WarningOrError(err, n, right == nil)
		if right == nil {
			right = util.NewStringLit("nil")
		}

		return util.NewBinaryExpr(left, operator, right), "bool",
			preStmts, postStmts, nil
	}

	// Convert "(0)" to "nil" when we are dealing with equality.
	if (operator == token.NEQ || operator == token.EQL) &&
		types.IsNullExpr(right) {
		t, err := types.ResolveType(p, leftType)
		ast.IsWarning(err, n)

		if t == "string" {
			p.AddImport("github.com/elliotchance/c2go/noarch")
			left = util.NewCallExpr("noarch.NullTerminatedString", left)
			right = &goast.BasicLit{
				Kind:  token.STRING,
				Value: `""`,
			}
		} else {
			right = goast.NewIdent("nil")
		}
	}

	if operator == token.ASSIGN {
		right, err = types.CastExpr(p, right, rightType, returnType)

		if ast.IsWarning(err, n) {
			right = util.NewStringLit("nil")
		}
	}

	return util.NewBinaryExpr(left, operator, right),
		types.ResolveTypeForBinaryOperator(p, n.Operator, leftType, rightType),
		preStmts, postStmts, nil
}
