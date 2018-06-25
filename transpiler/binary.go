// This file contains functions for transpiling binary operator expressions.

package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
)

// Comma problem. Example:
// for (int i=0,j=0;i+=1,j<5;i++,j++){...}
// For solving - we have to separate the
// binary operator "," to 2 parts:
// part 1(pre ): left part  - typically one or more some expessions
// part 2(stmt): right part - always only one expression, with or witout
//               logical operators like "==", "!=", ...
func transpileBinaryOperatorComma(n *ast.BinaryOperator, p *program.Program) (
	stmt goast.Stmt, preStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile operator comma : err = %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	left, err := transpileToStmts(n.Children()[0], p)
	if err != nil {
		return nil, nil, err
	}

	right, err := transpileToStmts(n.Children()[1], p)
	if err != nil {
		return nil, nil, err
	}

	if left == nil || right == nil {
		return nil, nil, fmt.Errorf("Cannot transpile binary operator comma: right = %v , left = %v", right, left)
	}

	preStmts = append(preStmts, left...)
	preStmts = append(preStmts, right...)

	if len(preStmts) >= 2 {
		return preStmts[len(preStmts)-1], preStmts[:len(preStmts)-1], nil
	}

	if len(preStmts) == 1 {
		return preStmts[0], nil, nil
	}
	return nil, nil, nil
}

func transpileBinaryOperator(n *ast.BinaryOperator, p *program.Program, exprIsStmt bool) (
	expr goast.Expr, eType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile BinaryOperator with type '%s' : result type = {%s}. Error: %v", n.Type, eType, err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	operator := getTokenForOperator(n.Operator)

	// Char overflow
	// BinaryOperator 0x2b74458 <line:506:7, col:18> 'int' '!='
	// |-ImplicitCastExpr 0x2b74440 <col:7, col:10> 'int' <IntegralCast>
	// | `-ImplicitCastExpr 0x2b74428 <col:7, col:10> 'char' <LValueToRValue>
	// |   `-...
	// `-ParenExpr 0x2b74408 <col:15, col:18> 'int'
	//   `-UnaryOperator 0x2b743e8 <col:16, col:17> 'int' prefix '-'
	//     `-IntegerLiteral 0x2b743c8 <col:17> 'int' 1
	if n.Operator == "!=" {
		var leftOk bool
		if l0, ok := n.ChildNodes[0].(*ast.ImplicitCastExpr); ok && l0.Type == "int" {
			if len(l0.ChildNodes) > 0 {
				if l1, ok := l0.ChildNodes[0].(*ast.ImplicitCastExpr); ok && l1.Type == "char" {
					leftOk = true
				}
			}
		}
		if leftOk {
			if r0, ok := n.ChildNodes[1].(*ast.ParenExpr); ok && r0.Type == "int" {
				if len(r0.ChildNodes) > 0 {
					if r1, ok := r0.ChildNodes[0].(*ast.UnaryOperator); ok && r1.IsPrefix && r1.Operator == "-" {
						if r2, ok := r1.ChildNodes[0].(*ast.IntegerLiteral); ok && r2.Type == "int" {
							r0.ChildNodes[0] = &ast.BinaryOperator{
								Type:     "int",
								Type2:    "int",
								Operator: "+",
								ChildNodes: []ast.Node{
									r1,
									&ast.IntegerLiteral{
										Type:  "int",
										Value: "256",
									},
								},
							}
						}
					}
				}
			}
		}
	}

	// Example of C code
	// a = b = 1
	// // Operation equal transpile from right to left
	// Solving:
	// b = 1, a = b
	// // Operation comma tranpile from left to right
	// If we have for example:
	// a = b = c = 1
	// then solution is:
	// c = 1, b = c, a = b
	// |-----------|
	// this part, created in according to
	// recursive working
	// Example of AST tree for problem:
	// |-BinaryOperator 0x2f17870 <line:13:2, col:10> 'int' '='
	// | |-DeclRefExpr 0x2f177d8 <col:2> 'int' lvalue Var 0x2f176d8 'x' 'int'
	// | `-BinaryOperator 0x2f17848 <col:6, col:10> 'int' '='
	// |   |-DeclRefExpr 0x2f17800 <col:6> 'int' lvalue Var 0x2f17748 'y' 'int'
	// |   `-IntegerLiteral 0x2f17828 <col:10> 'int' 1
	//
	// Example of AST tree for solution:
	// |-BinaryOperator 0x368e8d8 <line:13:2, col:13> 'int' ','
	// | |-BinaryOperator 0x368e820 <col:2, col:6> 'int' '='
	// | | |-DeclRefExpr 0x368e7d8 <col:2> 'int' lvalue Var 0x368e748 'y' 'int'
	// | | `-IntegerLiteral 0x368e800 <col:6> 'int' 1
	// | `-BinaryOperator 0x368e8b0 <col:9, col:13> 'int' '='
	// |   |-DeclRefExpr 0x368e848 <col:9> 'int' lvalue Var 0x368e6d8 'x' 'int'
	// |   `-ImplicitCastExpr 0x368e898 <col:13> 'int' <LValueToRValue>
	// |     `-DeclRefExpr 0x368e870 <col:13> 'int' lvalue Var 0x368e748 'y' 'int'
	if getTokenForOperator(n.Operator) == token.ASSIGN {
		switch c := n.Children()[1].(type) {
		case *ast.BinaryOperator:
			if getTokenForOperator(c.Operator) == token.ASSIGN {
				bSecond := ast.BinaryOperator{
					Type:     c.Type,
					Operator: "=",
				}
				bSecond.AddChild(n.Children()[0])

				var impl ast.ImplicitCastExpr
				impl.Type = c.Type
				impl.Kind = "LValueToRValue"
				impl.AddChild(c.Children()[0])
				bSecond.AddChild(&impl)

				var bComma ast.BinaryOperator
				bComma.Operator = ","
				bComma.Type = c.Type
				bComma.AddChild(c)
				bComma.AddChild(&bSecond)

				// goast.NewBinaryExpr takes care to wrap any AST children safely in a closure, if needed.
				return transpileBinaryOperator(&bComma, p, exprIsStmt)
			}
		}
	}

	// Example of C code
	// a = 1, b = a
	// Solving
	// a = 1; // preStmts
	// b = a; // n
	// Example of AST tree for problem:
	// |-BinaryOperator 0x368e8d8 <line:13:2, col:13> 'int' ','
	// | |-BinaryOperator 0x368e820 <col:2, col:6> 'int' '='
	// | | |-DeclRefExpr 0x368e7d8 <col:2> 'int' lvalue Var 0x368e748 'y' 'int'
	// | | `-IntegerLiteral 0x368e800 <col:6> 'int' 1
	// | `-BinaryOperator 0x368e8b0 <col:9, col:13> 'int' '='
	// |   |-DeclRefExpr 0x368e848 <col:9> 'int' lvalue Var 0x368e6d8 'x' 'int'
	// |   `-ImplicitCastExpr 0x368e898 <col:13> 'int' <LValueToRValue>
	// |     `-DeclRefExpr 0x368e870 <col:13> 'int' lvalue Var 0x368e748 'y' 'int'
	//
	// Example of AST tree for solution:
	// |-BinaryOperator 0x21a7820 <line:13:2, col:6> 'int' '='
	// | |-DeclRefExpr 0x21a77d8 <col:2> 'int' lvalue Var 0x21a7748 'y' 'int'
	// | `-IntegerLiteral 0x21a7800 <col:6> 'int' 1
	// |-BinaryOperator 0x21a78b0 <line:14:2, col:6> 'int' '='
	// | |-DeclRefExpr 0x21a7848 <col:2> 'int' lvalue Var 0x21a76d8 'x' 'int'
	// | `-ImplicitCastExpr 0x21a7898 <col:6> 'int' <LValueToRValue>
	// |   `-DeclRefExpr 0x21a7870 <col:6> 'int' lvalue Var 0x21a7748 'y' 'int'
	if getTokenForOperator(n.Operator) == token.COMMA {
		stmts, _, newPre, newPost, err := transpileToExpr(n.Children()[0], p, exprIsStmt)
		if err != nil {
			return nil, "unknown50", nil, nil, err
		}
		preStmts = append(preStmts, newPre...)
		preStmts = append(preStmts, util.NewExprStmt(stmts))
		preStmts = append(preStmts, newPost...)

		var st string
		stmts, st, newPre, newPost, err = transpileToExpr(n.Children()[1], p, exprIsStmt)
		if err != nil {
			return nil, "unknown51", nil, nil, err
		}
		// Theoretically , we don't have any preStmts or postStmts
		// from n.Children()[1]
		if len(newPre) > 0 || len(newPost) > 0 {
			p.AddMessage(p.GenerateWarningMessage(
				fmt.Errorf("Not support length pre or post stmts: {%d,%d}", len(newPre), len(newPost)), n))
		}
		return stmts, st, preStmts, postStmts, nil
	}

	left, leftType, newPre, newPost, err := atomicOperation(n.Children()[0], p)
	if err != nil {
		return nil, "unknown52", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	right, rightType, newPre, newPost, err := atomicOperation(n.Children()[1], p)
	if err != nil {
		return nil, "unknown53", nil, nil, err
	}
	var adjustPointerDiff int
	if types.IsPointer(p, leftType) && types.IsPointer(p, rightType) &&
		(operator == token.SUB ||
			operator == token.LSS || operator == token.GTR ||
			operator == token.LEQ || operator == token.GEQ) {
		baseSize, err := types.SizeOf(p, types.GetBaseType(leftType))
		if operator == token.SUB && err == nil && baseSize > 1 {
			adjustPointerDiff = baseSize
		}
		left, leftType, err = GetUintptrForPointer(p, left, leftType)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
		right, rightType, err = GetUintptrForPointer(p, right, rightType)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}
	if types.IsPointer(p, leftType) && types.IsPointer(p, rightType) &&
		(operator == token.EQL || operator == token.NEQ) &&
		leftType != "NullPointerType *" && rightType != "NullPointerType *" {
		left, leftType, err = GetUintptrForPointer(p, left, leftType)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
		right, rightType, err = GetUintptrForPointer(p, right, rightType)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	returnType := types.ResolveTypeForBinaryOperator(p, n.Operator, leftType, rightType)

	if operator == token.LAND || operator == token.LOR {
		left, err = types.CastExpr(p, left, leftType, "bool")
		p.AddMessage(p.GenerateWarningOrErrorMessage(err, n, left == nil))
		if left == nil {
			left = util.NewNil()
		}

		right, err = types.CastExpr(p, right, rightType, "bool")
		p.AddMessage(p.GenerateWarningOrErrorMessage(err, n, right == nil))
		if right == nil {
			right = util.NewNil()
		}

		resolvedLeftType, err := types.ResolveType(p, leftType)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}

		expr := util.NewBinaryExpr(left, operator, right, resolvedLeftType, exprIsStmt)

		return expr, "bool", preStmts, postStmts, nil
	}

	// The right hand argument of the shift left or shift right operators
	// in Go must be unsigned integers. In C, shifting with a negative shift
	// count is undefined behaviour (so we should be able to ignore that case).
	// To handle this, cast the shift count to a uint64.
	if operator == token.SHL || operator == token.SHR {
		right, err = types.CastExpr(p, right, rightType, "unsigned long long")
		p.AddMessage(p.GenerateWarningOrErrorMessage(err, n, right == nil))
		if right == nil {
			right = util.NewNil()
		}

		return util.NewBinaryExpr(left, operator, right, "uint64", exprIsStmt),
			leftType, preStmts, postStmts, nil
	}

	// pointer arithmetic
	if types.IsPointer(p, n.Type) {
		if operator == token.ADD || operator == token.SUB {
			if types.IsPointer(p, leftType) {
				expr, eType, newPre, newPost, err =
					pointerArithmetic(p, left, leftType, right, rightType, operator)
			} else {
				expr, eType, newPre, newPost, err =
					pointerArithmetic(p, right, rightType, left, leftType, operator)
			}
			if err != nil {
				return
			}
			if expr == nil {
				return nil, "", nil, nil, fmt.Errorf("Expr is nil")
			}
			preStmts, postStmts =
				combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

			return
		}
	}

	if operator == token.NEQ || operator == token.EQL ||
		operator == token.LSS || operator == token.GTR ||
		operator == token.LEQ || operator == token.GEQ ||
		operator == token.AND || operator == token.ADD ||
		operator == token.SUB || operator == token.MUL ||
		operator == token.QUO || operator == token.REM {

		// We may have to cast the right side to the same type as the left
		// side. This is a bit crude because we should make a better
		// decision of which type to cast to instead of only using the type
		// of the left side.
		if rightType != types.NullPointer {
			right, err = types.CastExpr(p, right, rightType, leftType)
			rightType = leftType
			p.AddMessage(p.GenerateWarningOrErrorMessage(err, n, right == nil))
		}

	}

	if operator == token.ASSIGN {
		// Memory allocation is translated into the Go-style.
		allocSize := getAllocationSizeNode(p, n.Children()[1])

		if allocSize != nil {
			right, newPre, newPost, err = generateAlloc(p, allocSize, leftType)
			if err != nil {
				p.AddMessage(p.GenerateWarningMessage(err, n))
				return nil, "", nil, nil, err
			}

			preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

		} else {
			right, err = types.CastExpr(p, right, rightType, returnType)

			if p.AddMessage(p.GenerateWarningMessage(err, n)) && right == nil {
				right = util.NewNil()
			}

		}
	}

	if operator == token.ADD_ASSIGN || operator == token.SUB_ASSIGN {
		right, err = types.CastExpr(p, right, rightType, returnType)
	}

	var resolvedLeftType = n.Type
	if !types.IsFunction(n.Type) && !types.IsTypedefFunction(p, n.Type) {
		resolvedLeftType, err = types.ResolveType(p, leftType)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}

	// Enum casting
	if operator != token.ASSIGN && strings.Contains(leftType, "enum") {
		left, err = types.CastExpr(p, left, leftType, "int")
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}

	// Enum casting
	if operator != token.ASSIGN && strings.Contains(rightType, "enum") {
		right, err = types.CastExpr(p, right, rightType, "int")
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}

	if left == nil {
		err = fmt.Errorf("left part of binary operation is nil. left : %#v", n.Children()[0])
		p.AddMessage(p.GenerateWarningMessage(err, n))
		return nil, "", nil, nil, err
	}

	if right == nil {
		err = fmt.Errorf("right part of binary operation is nil. right : %#v", n.Children()[1])
		p.AddMessage(p.GenerateWarningMessage(err, n))
		return nil, "", nil, nil, err
	}

	if adjustPointerDiff > 0 {
		expr := util.NewBinaryExpr(left, operator, right, resolvedLeftType, exprIsStmt)
		returnType = types.ResolveTypeForBinaryOperator(p, n.Operator, leftType, rightType)
		return util.NewBinaryExpr(expr, token.QUO, util.NewIntLit(adjustPointerDiff), returnType, exprIsStmt),
			returnType,
			preStmts, postStmts, nil
	}

	return util.NewBinaryExpr(left, operator, right, resolvedLeftType, exprIsStmt),
		types.ResolveTypeForBinaryOperator(p, n.Operator, leftType, rightType),
		preStmts, postStmts, nil
}

func foundCallExpr(n ast.Node) *ast.CallExpr {
	switch v := n.(type) {
	case *ast.ImplicitCastExpr, *ast.CStyleCastExpr:
		return foundCallExpr(n.Children()[0])
	case *ast.CallExpr:
		return v
	}
	return nil
}

// getAllocationSizeNode returns the node that, if evaluated, would return the
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
func getAllocationSizeNode(p *program.Program, node ast.Node) ast.Node {
	expr := foundCallExpr(node)

	if expr == nil || expr == (*ast.CallExpr)(nil) {
		return nil
	}

	functionName, _ := getNameOfFunctionFromCallExpr(p, expr)

	if functionName == "malloc" {
		// Is 1 always the body in this case? Might need to be more careful
		// to find the correct node.
		return expr.Children()[1]
	}

	if functionName == "calloc" {
		return &ast.BinaryOperator{
			Type:       "int",
			Operator:   "*",
			ChildNodes: expr.Children()[1:],
		}
	}

	// TODO: realloc() is not supported
	// https://github.com/elliotchance/c2go/issues/118
	//
	// Realloc will be treated as calloc which will almost certainly cause
	// bugs in your code.
	if functionName == "realloc" {
		return expr.Children()[2]
	}

	return nil
}

func generateAlloc(p *program.Program, allocSize ast.Node, leftType string) (
	right goast.Expr, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {

	allocSizeExpr, allocType, newPre, newPost, err := transpileToExpr(allocSize, p, false)

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	if err != nil {
		return nil, preStmts, postStmts, err
	}

	toType, err := types.ResolveType(p, leftType)
	if err != nil {
		return nil, preStmts, postStmts, err
	}
	allocSizeExpr, err = types.CastExpr(p, allocSizeExpr, allocType, "int")
	if err != nil {
		return nil, preStmts, postStmts, err
	}

	right = util.NewCallExpr(
		"noarch.Malloc",
		allocSizeExpr,
	)
	if toType != "unsafe.Pointer" {
		right = &goast.CallExpr{
			Fun: &goast.ParenExpr{
				X: util.NewTypeIdent(toType),
			},
			Args: []goast.Expr{right},
		}
	}
	return
}
