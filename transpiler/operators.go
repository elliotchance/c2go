// This file contains functions transpiling some general operator expressions.
// See binary.go and unary.go.

package transpiler

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"html/template"
	"strings"

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
	// Theoretically, length is must be zero
	if len(newPre) > 0 || len(newPost) > 0 {
		p.AddMessage(p.GenerateWarningMessage(
			fmt.Errorf("length of pre or post in body must be zero. {%d,%d}", len(newPre), len(newPost)), n))
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

	if !types.IsFunction(exprType) && exprType != "void" &&
		exprType != types.ToVoid {
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
// Example C code : ptr += i
// ptr = ((*int)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + (i)*unsafe.Sizeof(*ptr))))
// , where i  - left
//        '+' - operator
//      'ptr' - right
//      'int' - leftType transpiled in Go type
// Note:
// 1) rigthType MUST be 'int'
// 2) pointerArithmetic - implemented ONLY right part of formula
// 3) right is MUST be positive value, because impossible multiply uintptr to (-1)
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
	right, err = types.CastExpr(p, right, rightType, "int")
	if err != nil {
		return
	}
	if !types.IsPointer(p, leftType) {
		err = fmt.Errorf("left type is not a pointer : '%s'", leftType)
		return
	}

	resolvedLeftType, err := types.ResolveType(p, leftType)
	if err != nil {
		return
	}

	type pA struct {
		Name      string // name of variable: 'ptr'
		Type      string // type of variable: 'int','double'
		Condition string // condition : '-1' ,'(-1+2-2)'
		Operator  string // operator : '+', '-'
	}

	var s pA

	{
		var buf bytes.Buffer
		_ = printer.Fprint(&buf, token.NewFileSet(), left)
		s.Name = buf.String()
	}
	{
		var buf bytes.Buffer
		_ = printer.Fprint(&buf, token.NewFileSet(), right)
		s.Condition = buf.String()
	}
	s.Type = resolvedLeftType

	s.Operator = "+"
	if operator == token.SUB {
		s.Operator = "-"
	}

	var src string
	if util.IsAddressable(left) {
		src = `package main
func main(){
	a := (({{ .Type }})(unsafe.Pointer(uintptr(unsafe.Pointer({{ .Name }})) {{ .Operator }} (uintptr)({{ .Condition }})*unsafe.Sizeof(*{{ .Name }}))))
}`
	} else {
		src = `package main
func main(){
	a := (({{ .Type }})(func()unsafe.Pointer{
		tempVar := {{ .Name }}
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) {{ .Operator }} (uintptr)({{ .Condition }})*unsafe.Sizeof(*tempVar))
	}()))
}`
	}
	tmpl := template.Must(template.New("").Parse(src))
	var source bytes.Buffer
	err = tmpl.Execute(&source, s)
	if err != nil {
		err = fmt.Errorf("Cannot execute template. err = %v", err)
		return
	}

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	body := strings.Replace(source.String(), "&#43;", "+", -1)
	body = strings.Replace(body, "&#34;", "\"", -1)
	body = strings.Replace(body, "&#39;", "'", -1)
	body = strings.Replace(body, "&amp;", "&", -1)
	body = strings.Replace(body, "&lt;", "<", -1)
	body = strings.Replace(body, "&gt;", ">", -1)
	f, err := parser.ParseFile(fset, "", body, 0)
	if err != nil {
		err = fmt.Errorf("Cannot parse file. err = %v", err)
		return
	}

	p.AddImport("unsafe")

	return f.Decls[0].(*goast.FuncDecl).Body.List[0].(*goast.AssignStmt).Rhs[0],
		leftType, preStmts, postStmts, nil
}

func transpileCompoundAssignOperator(
	n *ast.CompoundAssignOperator, p *program.Program, exprIsStmt bool) (
	_ goast.Expr, _ string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileCompoundAssignOperator. err = %v", err)
		}
	}()

	operator := getTokenForOperator(n.Opcode)

	right, rightType, newPre, newPost, err := atomicOperation(n.Children()[1], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	left, leftType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// Pointer arithmetic
	if types.IsPointer(p, n.Type) &&
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
			X:  left,
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

	switch operator {
	case token.AND_ASSIGN, token.OR_ASSIGN, token.XOR_ASSIGN, token.AND_NOT_ASSIGN:
		right, err = types.CastExpr(p, right, rightType, leftType)
		p.AddMessage(p.GenerateWarningMessage(err, n))
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

	right, err = types.CastExpr(p, right, rightType, leftType)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
	}

	return util.NewBinaryExpr(left, operator, right, resolvedLeftType, exprIsStmt),
		n.Type, preStmts, postStmts, nil
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

func findUnaryWithInteger(node ast.Node) (*ast.UnaryOperator, bool) {
	switch n := node.(type) {
	case *ast.UnaryOperator:
		return n, true
	case *ast.ParenExpr:
		return findUnaryWithInteger(n.Children()[0])
	}
	return nil, false
}

func atomicOperation(n ast.Node, p *program.Program) (
	expr goast.Expr, exprType string, preStmts, postStmts []goast.Stmt, err error) {

	expr, exprType, preStmts, postStmts, err = transpileToExpr(n, p, false)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot create atomicOperation |%T|. err = %v", n, err)
		}
	}()

	switch v := n.(type) {
	case *ast.UnaryOperator:
		switch v.Operator {
		case "&", "*", "!", "-", "~":
			return
		}
		// UnaryOperator 0x252d798 <col:17, col:18> 'double' prefix '-'
		// `-FloatingLiteral 0x252d778 <col:18> 'double' 0.000000e+00
		if _, ok := v.Children()[0].(*ast.IntegerLiteral); ok {
			return
		}
		if _, ok := v.Children()[0].(*ast.FloatingLiteral); ok {
			return
		}

		// ++, -- anonymous functions are handled here below
		expr, exprType, preStmts, postStmts, err = transpileToExpr(n, p, true)

		// UnaryOperator 0x3001768 <col:204, col:206> 'int' prefix '++'
		// `-DeclRefExpr 0x3001740 <col:206> 'int' lvalue Var 0x303e888 'current_test' 'int'
		// OR
		// UnaryOperator 0x3001768 <col:204, col:206> 'int' postfix '++'
		// `-DeclRefExpr 0x3001740 <col:206> 'int' lvalue Var 0x303e888 'current_test' 'int'
		var varName string
		var vv *ast.DeclRefExpr
		if vv, err = getSoleChildDeclRefExpr(v); err == nil {
			varName = vv.Name

			var exprResolveType string
			exprResolveType, err = types.ResolveType(p, v.Type)
			if err != nil {
				return
			}

			// operators: ++, --
			if v.IsPrefix {
				// Example:
				// UnaryOperator 0x3001768 <col:204, col:206> 'int' prefix '++'
				// `-DeclRefExpr 0x3001740 <col:206> 'int' lvalue Var 0x303e888 'current_test' 'int'
				expr = util.NewAnonymousFunction(append(preStmts, &goast.ExprStmt{expr}),
					nil,
					util.NewIdent(varName),
					exprResolveType)
				preStmts = nil
				break
			}
			// Example:
			// UnaryOperator 0x3001768 <col:204, col:206> 'int' postfix '++'
			// `-DeclRefExpr 0x3001740 <col:206> 'int' lvalue Var 0x303e888 'current_test' 'int'
			expr = util.NewAnonymousFunction(preStmts,
				[]goast.Stmt{&goast.ExprStmt{expr}},
				util.NewIdent(varName),
				exprResolveType)
			preStmts = nil

			break
		}

		// UnaryOperator 0x358d470 <col:28, col:40> 'int' postfix '++'
		// `-MemberExpr 0x358d438 <col:28, col:36> 'int' lvalue .pos 0x358b538
		//   `-ArraySubscriptExpr 0x358d410 <col:28, col:34> 'struct struct_I_A':'struct struct_I_A' lvalue
		//     |-ImplicitCastExpr 0x358d3f8 <col:28> 'struct struct_I_A *' <ArrayToPointerDecay>
		//     | `-DeclRefExpr 0x358d3b0 <col:28> 'struct struct_I_A [2]' lvalue Var 0x358b6e8 'siia' 'struct struct_I_A [2]'
		//     `-IntegerLiteral 0x358d3d8 <col:33> 'int' 0
		varName = "tempVar"

		expr, exprType, preStmts, postStmts, err = transpileToExpr(v.Children()[0], p, false)
		if err != nil {
			return
		}

		body := append(preStmts, &goast.AssignStmt{
			Lhs: []goast.Expr{util.NewIdent(varName)},
			Tok: token.DEFINE,
			Rhs: []goast.Expr{&goast.UnaryExpr{
				Op: token.AND,
				X:  expr,
			}},
		})

		deferBody := postStmts
		postStmts = nil
		preStmts = nil

		switch v.Operator {
		case "++":
			expr = &goast.BinaryExpr{
				X:  &goast.StarExpr{X: util.NewIdent(varName)},
				Op: token.ADD_ASSIGN,
				Y:  &goast.BasicLit{Kind: token.INT, Value: "1"},
			}
		case "--":
			expr = &goast.BinaryExpr{
				X:  &goast.StarExpr{X: util.NewIdent(varName)},
				Op: token.SUB_ASSIGN,
				Y:  &goast.BasicLit{Kind: token.INT, Value: "1"},
			}
		}

		body = append(body, preStmts...)
		deferBody = append(deferBody, postStmts...)

		var exprResolveType string
		exprResolveType, err = types.ResolveType(p, v.Type)
		if err != nil {
			return
		}

		// operators: ++, --
		if v.IsPrefix {
			// Example:
			// UnaryOperator 0x3001768 <col:204, col:206> 'int' prefix '++'
			// `-DeclRefExpr 0x3001740 <col:206> 'int' lvalue Var 0x303e888 'current_test' 'int'
			expr = util.NewAnonymousFunction(append(body, &goast.ExprStmt{expr}), deferBody,
				&goast.StarExpr{
					X: util.NewIdent(varName),
				},
				exprResolveType)
			preStmts = nil
			postStmts = nil
			break
		}
		// Example:
		// UnaryOperator 0x3001768 <col:204, col:206> 'int' postfix '++'
		// `-DeclRefExpr 0x3001740 <col:206> 'int' lvalue Var 0x303e888 'current_test' 'int'
		expr = util.NewAnonymousFunction(body, append(deferBody, &goast.ExprStmt{expr}),
			&goast.StarExpr{
				X: util.NewIdent(varName),
			},
			exprResolveType)
		preStmts = nil
		postStmts = nil

	case *ast.CompoundAssignOperator:
		// CompoundAssignOperator 0x32911c0 <col:18, col:28> 'int' '-=' ComputeLHSTy='int' ComputeResultTy='int'
		// |-DeclRefExpr 0x3291178 <col:18> 'int' lvalue Var 0x328df60 'iterator' 'int'
		// `-IntegerLiteral 0x32911a0 <col:28> 'int' 2
		if vv, ok := v.Children()[0].(*ast.DeclRefExpr); ok {
			var varName string
			varName = vv.Name

			var exprResolveType string
			exprResolveType, err = types.ResolveType(p, v.Type)
			if err != nil {
				return
			}

			// since we will explicitly use an anonymous function, we can transpileToExpr as a statement
			expr, exprType, preStmts, postStmts, err = transpileToExpr(n, p, true)
			expr = util.NewAnonymousFunction(append(preStmts, &goast.ExprStmt{expr}),
				postStmts,
				util.NewIdent(varName),
				exprResolveType)
			preStmts = nil
			postStmts = nil
			break
		}
		// CompoundAssignOperator 0x27906c8 <line:450:2, col:6> 'double' '+=' ComputeLHSTy='double' ComputeResultTy='double'
		// |-UnaryOperator 0x2790670 <col:2, col:3> 'double' lvalue prefix '*'
		// | `-ImplicitCastExpr 0x2790658 <col:3> 'double *' <LValueToRValue>
		// |   `-DeclRefExpr 0x2790630 <col:3> 'double *' lvalue Var 0x2790570 'p' 'double *'
		// `-IntegerLiteral 0x32911a0 <col:28> 'int' 2
		if vv, ok := v.Children()[0].(*ast.UnaryOperator); ok && vv.IsPrefix && vv.Operator == "*" {
			if vvv, ok := vv.Children()[0].(*ast.ImplicitCastExpr); ok {
				if vvvv, ok := vvv.Children()[0].(*ast.DeclRefExpr); ok {
					if types.IsPointer(p, vvvv.Type) {
						var varName string
						varName = vvvv.Name

						var exprResolveType string
						exprResolveType, err = types.ResolveType(p, v.Type)
						if err != nil {
							return
						}

						expr = util.NewAnonymousFunction(append(preStmts, &goast.ExprStmt{expr}),
							postStmts,
							&goast.UnaryExpr{
								Op: token.AND,
								X:  util.NewIdent(varName),
							},
							exprResolveType)
						preStmts = nil
						postStmts = nil
						break
					}
				}
			}
		}

		// CompoundAssignOperator 0x32911c0 <col:18, col:28> 'int' '-=' ComputeLHSTy='int' ComputeResultTy='int'
		// |-DeclRefExpr 0x3291178 <col:18> 'int' lvalue Var 0x328df60 'iterator' 'int'
		// `-IntegerLiteral 0x32911a0 <col:28> 'int' 2
		varName := "tempVar"
		expr, exprType, preStmts, postStmts, err = transpileToExpr(v.Children()[0], p, false)
		if err != nil {
			return
		}
		body := append(preStmts, &goast.AssignStmt{
			Lhs: []goast.Expr{util.NewIdent(varName)},
			Tok: token.DEFINE,
			Rhs: []goast.Expr{&goast.UnaryExpr{
				Op: token.AND,
				X:  expr,
			}},
		})
		preStmts = nil

		// CompoundAssignOperator 0x27906c8 <line:450:2, col:6> 'double' '+=' ComputeLHSTy='double' ComputeResultTy='double'
		// |-UnaryOperator 0x2790670 <col:2, col:3> 'double' lvalue prefix '*'
		// | `-ImplicitCastExpr 0x2790658 <col:3> 'double *' <LValueToRValue>
		// |   `-DeclRefExpr 0x2790630 <col:3> 'double *' lvalue Var 0x2790570 'p' 'double *'
		// `-ImplicitCastExpr 0x27906b0 <col:6> 'double' <IntegralToFloating>
		//   `-IntegerLiteral 0x2790690 <col:6> 'int' 1
		var newPre, newPost []goast.Stmt
		expr, exprType, newPre, newPost, err = atomicOperation(v.Children()[1], p)
		if err != nil {
			return
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

		var exprResolveType string
		exprResolveType, err = types.ResolveType(p, v.Type)
		if err != nil {
			return
		}

		body = append(preStmts, body...)
		body = append(body, &goast.AssignStmt{
			Lhs: []goast.Expr{&goast.StarExpr{
				X: util.NewIdent(varName),
			}},
			Tok: getTokenForOperator(v.Opcode),
			Rhs: []goast.Expr{expr},
		})

		expr = util.NewAnonymousFunction(body, postStmts,
			&goast.StarExpr{
				X: util.NewIdent(varName),
			},
			exprResolveType)
		preStmts = nil
		postStmts = nil

	case *ast.ParenExpr:
		// ParenExpr 0x3c42468 <col:18, col:40> 'int'
		return atomicOperation(v.Children()[0], p)

	case *ast.ImplicitCastExpr:
		if _, ok := v.Children()[0].(*ast.MemberExpr); ok {
			return
		}
		if _, ok := v.Children()[0].(*ast.IntegerLiteral); ok {
			return
		}

		// for case : overflow char
		// ImplicitCastExpr 0x2027358 <col:6, col:7> 'char' <IntegralCast>
		// `-UnaryOperator 0x2027338 <col:6, col:7> 'int' prefix '-'
		//   `-IntegerLiteral 0x2027318 <col:7> 'int' 1
		//
		// another example :
		// ImplicitCastExpr 0x2982630 <col:11, col:14> 'char' <IntegralCast>
		// `-ParenExpr 0x2982610 <col:11, col:14> 'int'
		//   `-UnaryOperator 0x29825f0 <col:12, col:13> 'int' prefix '-'
		//     `-IntegerLiteral 0x29825d0 <col:13> 'int' 1
		if v.Type == "char" {
			if len(v.Children()) == 1 {
				if u, ok := findUnaryWithInteger(n.Children()[0]); ok {
					if u.IsPrefix && u.Type == "int" && u.Operator == "-" {
						if _, ok := u.Children()[0].(*ast.IntegerLiteral); ok {
							return transpileToExpr(&ast.BinaryOperator{
								Type:     "int",
								Type2:    "int",
								Operator: "+",
								ChildNodes: []ast.Node{
									u,
									&ast.IntegerLiteral{
										Type:  "int",
										Value: "256",
									},
								},
							}, p, false)
						}
					}
				}
			}
		}

		expr, exprType, preStmts, postStmts, err = atomicOperation(v.Children()[0], p)
		if err != nil {
			return nil, "", nil, nil, err
		}
		if exprType == types.NullPointer {
			return
		}
		if !types.IsFunction(exprType) && !strings.ContainsAny(v.Type, "[]") {
			expr, err = types.CastExpr(p, expr, exprType, v.Type)
			if err != nil {
				return nil, "", nil, nil, err
			}
			exprType = v.Type
		}
		return

	case *ast.BinaryOperator:
		switch v.Operator {
		case ",":
			// BinaryOperator 0x35b95e8 <col:29, col:51> 'int' ','
			// |-UnaryOperator 0x35b94b0 <col:29, col:31> 'int' postfix '++'
			// | `-DeclRefExpr 0x35b9488 <col:29> 'int' lvalue Var 0x35b8dc8 't' 'int'
			// `-CompoundAssignOperator 0x35b95b0 <col:36, col:51> 'int' '+=' ComputeLHSTy='int' ComputeResultTy='int'
			//   |-MemberExpr 0x35b9558 <col:36, col:44> 'int' lvalue .pos 0x35b8730
			//   | `-ArraySubscriptExpr 0x35b9530 <col:36, col:42> 'struct struct_I_A4':'struct struct_I_A4' lvalue
			//   |   |-ImplicitCastExpr 0x35b9518 <col:36> 'struct struct_I_A4 *' <ArrayToPointerDecay>
			//   |   | `-DeclRefExpr 0x35b94d0 <col:36> 'struct struct_I_A4 [2]' lvalue Var 0x35b88d8 'siia' 'struct struct_I_A4 [2]'
			//   |   `-IntegerLiteral 0x35b94f8 <col:41> 'int' 0
			//   `-IntegerLiteral 0x35b9590 <col:51> 'int' 1

			// `-BinaryOperator 0x3c42440 <col:19, col:32> 'int' ','
			//   |-BinaryOperator 0x3c423d8 <col:19, col:30> 'int' '='
			//   | |-DeclRefExpr 0x3c42390 <col:19> 'int' lvalue Var 0x3c3cf60 'iterator' 'int'
			//   | `-IntegerLiteral 0x3c423b8 <col:30> 'int' 0
			//   `-ImplicitCastExpr 0x3c42428 <col:32> 'int' <LValueToRValue>
			//     `-DeclRefExpr 0x3c42400 <col:32> 'int' lvalue Var 0x3c3cf60 'iterator' 'int'
			varName := "tempVar"

			expr, exprType, preStmts, postStmts, err = transpileToExpr(v.Children()[0], p, true)
			if err != nil {
				return
			}

			inBody := combineStmts(&goast.ExprStmt{expr}, preStmts, postStmts)
			preStmts = nil
			postStmts = nil

			expr, exprType, preStmts, postStmts, err = atomicOperation(v.Children()[1], p)
			if err != nil {
				return
			}

			if v, ok := expr.(*goast.CallExpr); ok {
				if vv, ok := v.Fun.(*goast.FuncLit); ok {
					vv.Body.List = append(inBody, vv.Body.List...)
					break
				}
			}

			body := append(inBody, preStmts...)
			preStmts = nil

			body = append(body, &goast.AssignStmt{
				Lhs: []goast.Expr{util.NewIdent(varName)},
				Tok: token.DEFINE,
				Rhs: []goast.Expr{&goast.UnaryExpr{
					Op: token.AND,
					X:  expr,
				}},
			})

			var exprResolveType string
			exprResolveType, err = types.ResolveType(p, v.Type)
			if err != nil {
				return
			}

			expr = util.NewAnonymousFunction(body, postStmts,
				&goast.UnaryExpr{
					Op: token.MUL,
					X:  util.NewIdent(varName),
				},
				exprResolveType)
			preStmts = nil
			postStmts = nil
			exprType = v.Type
			return

		case "=":
			// BinaryOperator 0x2a230c0 <col:8, col:13> 'int' '='
			// |-UnaryOperator 0x2a23080 <col:8, col:9> 'int' lvalue prefix '*'
			// | `-ImplicitCastExpr 0x2a23068 <col:9> 'int *' <LValueToRValue>
			// |   `-DeclRefExpr 0x2a23040 <col:9> 'int *' lvalue Var 0x2a22f20 'a' 'int *'
			// `-IntegerLiteral 0x2a230a0 <col:13> 'int' 42

			// VarDecl 0x328dc50 <col:3, col:29> col:13 used d 'int' cinit
			// `-BinaryOperator 0x328dd98 <col:17, col:29> 'int' '='
			//   |-DeclRefExpr 0x328dcb0 <col:17> 'int' lvalue Var 0x328dae8 'a' 'int'
			//   `-BinaryOperator 0x328dd70 <col:21, col:29> 'int' '='
			//     |-DeclRefExpr 0x328dcd8 <col:21> 'int' lvalue Var 0x328db60 'b' 'int'
			//     `-BinaryOperator 0x328dd48 <col:25, col:29> 'int' '='
			//       |-DeclRefExpr 0x328dd00 <col:25> 'int' lvalue Var 0x328dbd8 'c' 'int'
			//       `-IntegerLiteral 0x328dd28 <col:29> 'int' 42

			var body []goast.Stmt
			varName := "tempVar"

			var exprResolveType string
			exprResolveType, err = types.ResolveType(p, v.Type)
			if err != nil {
				return
			}

			e, _, newPre, newPost, _ := transpileToExpr(v, p, true)
			if assign, ok := e.(*goast.BinaryExpr); !ok || assign.Op != token.ASSIGN {
				panic("not a valid assignment")
			} else {
				body = append(body, &goast.AssignStmt{
					Lhs: []goast.Expr{util.NewIdent(varName)},
					Tok: token.DEFINE,
					Rhs: []goast.Expr{assign.Y},
				})
				body = append(body, &goast.ExprStmt{
					&goast.BinaryExpr{
						X:  assign.X,
						Op: token.ASSIGN,
						Y:  util.NewIdent(varName),
					},
				})
			}

			body = combineMultipleStmts(body, newPre, newPost)

			preStmts = nil
			postStmts = nil

			var returnValue goast.Expr = util.NewIdent(varName)

			expr = util.NewAnonymousFunction(body,
				nil,
				returnValue,
				exprResolveType)
			expr = &goast.ParenExpr{
				X:      expr,
				Lparen: 1,
			}
		}

	}

	return
}

// getDeclRefExpr - find ast DeclRefExpr
// Examples of input ast trees:
// UnaryOperator 0x2a23080 <col:8, col:9> 'int' lvalue prefix '*'
// `-ImplicitCastExpr 0x2a23068 <col:9> 'int *' <LValueToRValue>
//   `-DeclRefExpr 0x2a23040 <col:9> 'int *' lvalue Var 0x2a22f20 'a' 'int *'
//
// DeclRefExpr 0x328dd00 <col:25> 'int' lvalue Var 0x328dbd8 'c' 'int'
func getDeclRefExpr(n ast.Node) (*ast.DeclRefExpr, bool) {
	switch v := n.(type) {
	case *ast.DeclRefExpr:
		return v, true
	case *ast.ImplicitCastExpr:
		return getDeclRefExpr(n.Children()[0])
	case *ast.UnaryOperator:
		return getDeclRefExpr(n.Children()[0])
	}
	return nil, false
}
