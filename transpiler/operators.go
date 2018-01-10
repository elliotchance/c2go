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
	if expr == nil {
		err = fmt.Errorf("Expr is nil")
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

// pointerArithmetic - operations between 'int' and pointer
// Example:
// code : ptr += 1
// AST:
// CompoundAssignOperator 0x2738b20 <line:300:3, col:10> 'int *' '+=' ComputeLHSTy='int *' ComputeResultTy='int *'
// |-DeclRefExpr 0x2738ad8 <col:3> 'int *' lvalue Var 0x2737a90 'ptr' 'int *'
// `-IntegerLiteral 0x2738b00 <col:10> 'int' 1
// Solution on meta language:
// ptr = func() { return ptr + 1 }()
// Example of solution on Go:
// intArray := [...]int{1, 2}
// intPtr := &intArray[0] // type of intPtr is '* int' on Go
// intPtr = (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(intPtr)) + unsafe.Sizeof(intArray[0])))
// for our case :
// ptr = ([]int)(unsafe.Pointer(uintptr(unsafe.Pointer(&ptr[0])) + i * unsafe.Sizeof(ptr[0])))
// , where i  - left
//        '+' - operator
//      'ptr' - right
//      'int' - leftType transpiled in Go type
// Note:
// rigthType MUST be 'int'
// pointerArithmetic - implemented ONLY right part of formula
func pointerArithmetic(p *program.Program,
	left goast.Expr, leftType string,
	right goast.Expr, rightType string,
	operator token.Token) (
	_ goast.Expr, _ string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile pointerArithmetic. err = %v", err)
		}
	}()
	if rightType != "int" {
		err = fmt.Errorf("right type is not 'int' : '%s'", rightType)
		return
	}
	if !types.IsPointer(leftType) {
		err = fmt.Errorf("left type is not a pointer : '%s'", leftType)
		return
	}

	resolvedLeftType, err := types.ResolveType(p, leftType)
	if err != nil {
		return
	}

	expr := &goast.StarExpr{
		Star: 1,
		X: &goast.CallExpr{
			Fun: &goast.ParenExpr{
				Lparen: 1,
				X: &goast.StarExpr{
					Star: 1,
					X: &goast.ParenExpr{
						Lparen: 1,
						X:      goast.NewIdent(resolvedLeftType), // Type
					},
				},
			},
			Lparen: 1,
			Args: []goast.Expr{
				&goast.CallExpr{
					Fun: &goast.SelectorExpr{
						X:   goast.NewIdent("unsafe"),
						Sel: goast.NewIdent("Pointer"),
					},
					Lparen: 1,
					Args: []goast.Expr{
						&goast.BinaryExpr{
							X: &goast.CallExpr{
								Fun:    goast.NewIdent("uintptr"),
								Lparen: 1,
								Args: []goast.Expr{
									&goast.CallExpr{
										Fun: &goast.SelectorExpr{
											X:   goast.NewIdent("unsafe"),
											Sel: goast.NewIdent("Pointer"),
										},
										Lparen: 1,
										Args: []goast.Expr{
											&goast.UnaryExpr{
												Op: token.AND, // &
												X: &goast.IndexExpr{
													X:      left, // ptr
													Lbrack: 1,
													Index: &goast.BasicLit{
														Kind:  token.INT,
														Value: "0",
													},
												},
											},
										},
									},
								},
							},
							Op: token.ADD, // operation
							Y: &goast.BinaryExpr{
								X:  &goast.ParenExpr{Lparen: 1, X: right}, // i
								Op: token.MUL,                             // *
								Y: &goast.CallExpr{
									Fun: &goast.SelectorExpr{
										X:   goast.NewIdent("unsafe"),
										Sel: goast.NewIdent("Sizeof"),
									},
									Lparen: 1,
									Args: []goast.Expr{
										&goast.IndexExpr{
											X:      left, // ptr
											Lbrack: 1,
											Index: &goast.BasicLit{
												Kind:  token.INT,
												Value: "0",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	p.AddImport("unsafe")

	return expr, leftType, preStmts, postStmts, nil
}

func transpileCompoundAssignOperator(n *ast.CompoundAssignOperator, p *program.Program, exprIsStmt bool) (
	_ goast.Expr, _ string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileCompoundAssignOperator. err = %v", err)
		}
	}()

	operator := getTokenForOperator(n.Opcode)

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

	// Pointer arithmetic
	if types.IsPointer(n.Type) &&
		(operator == token.ADD_ASSIGN || operator == token.SUB_ASSIGN) {
		operator = convertToWithoutAssign(operator)
		v, vType, newPre, newPost, err := pointerArithmetic(p, left, leftType, right, rightType, operator)
		if err != nil {
			return nil, "", nil, nil, err
		}
		if v == nil {
			return nil, "", nil, nil, fmt.Errorf("Expr is nil")
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		v = &goast.BinaryExpr{
			X:  util.NewIdent(n.Children()[0].(*ast.DeclRefExpr).Name),
			Op: token.ASSIGN,
			Y:  v,
		}
		return v, vType, preStmts, postStmts, nil
	}

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

	if right == nil {
		err = fmt.Errorf("Right part is nil. err = %v", err)
		return nil, "", nil, nil, err
	}
	if left == nil {
		err = fmt.Errorf("Left part is nil. err = %v", err)
		return nil, "", nil, nil, err
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

func convertToWithoutAssign(operator token.Token) token.Token {
	switch operator {
	case token.ADD_ASSIGN: // "+="
		return token.ADD
	case token.SUB_ASSIGN: // "-="
		return token.SUB
	case token.MUL_ASSIGN: // "*="
		return token.MUL
	case token.QUO_ASSIGN: // "/="
		return token.QUO
	}
	panic(fmt.Sprintf("not support operator: %v", operator))
}
