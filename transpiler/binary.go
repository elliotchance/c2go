// This file contains functions for transpiling binary operator expressions.

package transpiler

import (
	goast "go/ast"
	"go/token"
	"reflect"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/traverse"
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
		right = goast.NewIdent("nil")
	}

	if operator == token.ASSIGN {
		// Memory allocation is translated into the Go-style.
		allocSize := GetAllocationSizeNode(n.Children[1])

		if allocSize != nil {
			allocSizeExpr, _, newPre, newPost, err := transpileToExpr(allocSize, p)
			preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

			if err != nil {
				return nil, "", preStmts, postStmts, err
			}

			derefType, err := types.GetDereferenceType(leftType)
			if err != nil {
				return nil, "", preStmts, postStmts, err
			}

			toType, err := types.ResolveType(p, leftType)
			if err != nil {
				return nil, "", preStmts, postStmts, err
			}

			elementSize, err := types.SizeOf(p, derefType)
			if err != nil {
				return nil, "", preStmts, postStmts, err
			}

			right = util.NewCallExpr(
				"make",
				util.NewStringLit(toType),
				util.NewBinaryExpr(allocSizeExpr, token.QUO, util.NewIntLit(elementSize)),
			)
		} else {
			right, err = types.CastExpr(p, right, rightType, returnType)

			if ast.IsWarning(err, n) && right == nil {
				right = util.NewStringLit("nil")
			}
		}
	}

	return util.NewBinaryExpr(left, operator, right),
		types.ResolveTypeForBinaryOperator(p, n.Operator, leftType, rightType),
		preStmts, postStmts, nil
}

// GetAllocationSizeNode returns the node that, if evaluated, would return the
// size (in bytes) of a memory allocation operation. For example:
//
//     (int *)malloc(sizeof(int))
//
// Would return the node that represents the "sizeof(int)".
//
// If the node does not represent an allocation operation (such as calling
// malloc, calloc, realloc, etc.) then nil is returned.
func GetAllocationSizeNode(node ast.Node) ast.Node {
	exprs := traverse.GetAllNodesOfType(node,
		reflect.TypeOf((*ast.CallExpr)(nil)))

	for _, expr := range exprs {
		functionName, _ := getNameOfFunctionFromCallExpr(expr.(*ast.CallExpr))

		if functionName == "malloc" ||
			functionName == "calloc" ||
			functionName == "realloc" {
			// Is 1 always the body in this case? Might need to be more careful
			// to find the correct node.
			return expr.(*ast.CallExpr).Children[1]
		}
	}

	return nil
}
