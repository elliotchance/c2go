// This file contains functions for transpiling binary operator expressions.

package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"reflect"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/traverse"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
)

func transpileBinaryOperator(n *ast.BinaryOperator, p *program.Program) (
	goast.Expr, string, []goast.Stmt, []goast.Stmt, error) {
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

	if operator == token.LAND || operator == token.LOR {
		left, err = types.CastExpr(p, left, leftType, "bool")
		p.AddMessage(ast.GenerateWarningOrErrorMessage(err, n, left == nil))
		if left == nil {
			left = util.NewNil()
		}

		right, err = types.CastExpr(p, right, rightType, "bool")
		p.AddMessage(ast.GenerateWarningOrErrorMessage(err, n, right == nil))
		if right == nil {
			right = util.NewNil()
		}

		return util.NewBinaryExpr(left, operator, right), "bool",
			preStmts, postStmts, nil
	}

	// The right hand argument of the shift left or shift right operators
	// in Go must be unsigned integers. In C, shifting with a negative shift
	// count is undefined behaviour (so we should be able to ignore that case).
	// To handle this, cast the shift count to a uint64.
	if operator == token.SHL || operator == token.SHR {
		right, err = types.CastExpr(p, right, rightType, "unsigned long long")
		p.AddMessage(ast.GenerateWarningOrErrorMessage(err, n, right == nil))
		if right == nil {
			right = util.NewNil()
		}

		return util.NewBinaryExpr(left, operator, right), leftType,
			preStmts, postStmts, nil
	}

	if operator == token.NEQ || operator == token.EQL {
		// Convert "(0)" to "nil" when we are dealing with equality.
		if types.IsNullExpr(right) {
			right = util.NewNil()
		} else {
			// We may have to cast the right side to the same type as the left
			// side. This is a bit crude because we should make a better
			// decision of which type to cast to instead of only using the type
			// of the left side.
			right, err = types.CastExpr(p, right, rightType, leftType)
			p.AddMessage(ast.GenerateWarningOrErrorMessage(err, n, right == nil))
		}
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
				util.NewTypeIdent(toType),
				util.NewBinaryExpr(allocSizeExpr, token.QUO, util.NewIntLit(elementSize)),
			)
		} else {
			right, err = types.CastExpr(p, right, rightType, returnType)

			if _, ok := right.(*goast.UnaryExpr); ok {
				deref, err := types.GetDereferenceType(rightType)

				if !p.AddMessage(ast.GenerateWarningMessage(err, n)) {
					// This is some hackey to convert a reference to a variable
					// into a slice that points to the same location. It will
					// look similar to:
					//
					//     (*[1]int)(unsafe.Pointer(&a))[:]
					//
					p.AddImport("unsafe")
					right = &goast.SliceExpr{
						X: util.NewCallExpr(
							fmt.Sprintf("(*[1]%s)", deref),
							util.NewCallExpr("unsafe.Pointer", right),
						),
					}
				}
			}

			if p.AddMessage(ast.GenerateWarningMessage(err, n)) && right == nil {
				right = util.NewNil()
			}

			// Construct code for assigning value to an union field
			member_expr, ok := n.Children[0].(*ast.MemberExpr)
			if ok {
				ref := member_expr.GetDeclRef()
				if ref != nil {
					typename, err := types.ResolveType(p, ref.Type)
					if err != nil {
						return nil, "", preStmts, postStmts, err
					}

					if typename[0] == '*' {
						typename = typename[1:]
					}

					union := p.GetStruct(typename)
					if union.IsUnion {
						resExpr := &goast.CallExpr{
							Fun: &goast.SelectorExpr{
								X:   goast.NewIdent(ref.Name),
								Sel: goast.NewIdent("Set" + strings.Title(member_expr.Name)),
							},
							Args: []goast.Expr{
								right,
							},
						}

						resType := types.ResolveTypeForBinaryOperator(p, n.Operator, leftType, rightType)

						return resExpr, resType, preStmts, postStmts, nil
					}
				}
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
//
// In the case of calloc() it will return a new BinaryExpr that multiplies both
// arguments.
func GetAllocationSizeNode(node ast.Node) ast.Node {
	exprs := traverse.GetAllNodesOfType(node,
		reflect.TypeOf((*ast.CallExpr)(nil)))

	for _, expr := range exprs {
		functionName, _ := getNameOfFunctionFromCallExpr(expr.(*ast.CallExpr))

		if functionName == "malloc" {
			// Is 1 always the body in this case? Might need to be more careful
			// to find the correct node.
			return expr.(*ast.CallExpr).Children[1]
		}

		if functionName == "calloc" {
			return &ast.BinaryOperator{
				Type:     "int",
				Operator: "*",
				Children: expr.(*ast.CallExpr).Children[1:],
			}
		}

		// TODO: realloc() is not supported
		// https://github.com/elliotchance/c2go/issues/118
		//
		// Realloc will be treated as calloc which will almost certainly cause
		// bugs in your code.
		if functionName == "realloc" {
			return expr.(*ast.CallExpr).Children[2]
		}
	}

	return nil
}
