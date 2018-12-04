// This file contains functions for transpiling unary operator expressions.

package transpiler

import (
	"fmt"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"

	goast "go/ast"
	"go/token"
)

func transpileUnaryOperatorInc(n *ast.UnaryOperator, p *program.Program, operator token.Token, exprIsStmt bool) (
	expr goast.Expr, eType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileUnaryOperatorInc. err = %v", err)
		}
	}()

	if !(operator == token.INC || operator == token.DEC) {
		err = fmt.Errorf("not acceptable operator '%v'", operator)
		return
	}

	if types.IsPointer(p, n.Type) {
		switch operator {
		case token.INC:
			operator = token.ADD
		case token.DEC:
			operator = token.SUB
		}
		var e ast.Node
		e, err = getSoleChildIncrementable(n)
		if err != nil {
			return
		}

		var left goast.Expr
		var leftType string
		var newPre, newPost []goast.Stmt
		left, leftType, newPre, newPost, err = transpileToExpr(e, p, false)
		if err != nil {
			return
		}

		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

		rightType := "int"
		right := &goast.BasicLit{
			Kind:  token.INT,
			Value: "1",
		}

		expr, eType, newPre, newPost, err = pointerArithmetic(p, left, leftType, right, rightType, operator)
		if err != nil {
			return
		}
		if expr == nil {
			return nil, "", nil, nil, fmt.Errorf("Expr is nil")
		}

		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

		expr = &goast.BinaryExpr{
			X:  left,
			Op: token.ASSIGN,
			Y:  expr,
		}
		if !exprIsStmt {
			var lType string
			lType, err = types.ResolveType(p, leftType)
			if err != nil {
				return
			}
			expr = util.NewAnonymousFunction([]goast.Stmt{&goast.ExprStmt{
				X: expr,
			}}, nil, left, lType)
		}
		return
	}

	if v, ok := n.Children()[0].(*ast.DeclRefExpr); ok {
		switch n.Operator {
		case "++":
			return &goast.BinaryExpr{
				X:  util.NewIdent(v.Name),
				Op: token.ADD_ASSIGN,
				Y:  &goast.BasicLit{Kind: token.INT, Value: "1"},
			}, n.Type, nil, nil, nil
		case "--":
			return &goast.BinaryExpr{
				X:  util.NewIdent(v.Name),
				Op: token.SUB_ASSIGN,
				Y:  &goast.BasicLit{Kind: token.INT, Value: "1"},
			}, n.Type, nil, nil, nil
		}
	}

	// Unfortunately we cannot use the Go increment operators because we are not
	// providing any position information for tokens. This means that the ++/--
	// would be placed before the expression and would be invalid in Go.
	//
	// Until it can be properly fixed (can we trick Go into to placing it after
	// the expression with a magic position?) we will have to return a
	// BinaryExpr with the same functionality.

	binaryOperator := "+="
	if operator == token.DEC {
		binaryOperator = "-="
	}

	return transpileBinaryOperator(&ast.BinaryOperator{
		Type:     n.Type,
		Operator: binaryOperator,
		ChildNodes: []ast.Node{
			n.Children()[0], &ast.IntegerLiteral{
				Type:       "int",
				Value:      "1",
				ChildNodes: []ast.Node{},
			},
		},
	}, p, exprIsStmt)
}

func getSoleChildIncrementable(n ast.Node) (result ast.Node, err error) {
	children := n.Children()
	if len(children) != 1 {
		return nil, fmt.Errorf("expected one child node, got %d", len(children))
	}
	switch c := children[0].(type) {
	case *ast.ParenExpr:
		return getSoleChildIncrementable(c)
	case *ast.DeclRefExpr, *ast.MemberExpr, *ast.UnaryOperator:
		return c, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", c)
	}
}

func getSoleChildDeclRefExpr(n *ast.UnaryOperator) (result *ast.DeclRefExpr, err error) {
	var inspect ast.Node = n
	for {
		if inspect == nil {
			break
		}
		if ret, ok := inspect.Children()[0].(*ast.DeclRefExpr); ok {
			return ret, nil
		}
		if len(inspect.Children()) > 1 {
			return nil, fmt.Errorf("node has to many children: %T", inspect)
		} else if len(inspect.Children()) == 1 {
			if _, ok := inspect.Children()[0].(*ast.ParenExpr); !ok {
				err = fmt.Errorf("unsupported type %T", n.Children()[0])
				return
			}
			inspect = inspect.Children()[0]
		} else {
			break
		}
	}
	return nil, fmt.Errorf("could not find supported type DeclRefExpr")
}

func transpileUnaryOperatorNot(n *ast.UnaryOperator, p *program.Program) (
	goast.Expr, string, []goast.Stmt, []goast.Stmt, error) {
	e, eType, preStmts, postStmts, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}
	// null in C is zero
	if eType == types.NullPointer {
		e = &goast.BasicLit{
			Kind:  token.INT,
			Value: "0",
		}
		eType = "int"
	}

	if eType == "bool" {
		return &goast.UnaryExpr{
			X:  e,
			Op: token.NOT,
		}, "bool", preStmts, postStmts, nil
	}

	if strings.HasSuffix(eType, "*") {
		// `!pointer` has to be converted to `pointer == nil`
		return &goast.BinaryExpr{
			X:  e,
			Op: token.EQL,
			Y:  util.NewIdent("nil"),
		}, "bool", preStmts, postStmts, nil
	}

	t, err := types.ResolveType(p, eType)
	p.AddMessage(p.GenerateWarningMessage(err, n))

	if t == "*byte" {
		return util.NewUnaryExpr(
			token.NOT, util.NewCallExpr("noarch.CStringIsNull", e),
		), "bool", preStmts, postStmts, nil
	}

	// only if added "stdbool.h"
	if p.IncludeHeaderIsExists("stdbool.h") {
		if t == "_Bool" {
			t = "int8"
			e = &goast.CallExpr{
				Fun:    goast.NewIdent("int8"),
				Lparen: 1,
				Args:   []goast.Expr{e},
			}
		}
	}

	p.AddImport("github.com/elliotchance/c2go/noarch")

	functionName := fmt.Sprintf("noarch.Not%s",
		util.GetExportedName(t))

	return util.NewCallExpr(functionName, e),
		eType, preStmts, postStmts, nil
}

// tranpileUnaryOperatorAmpersant - operator ampersant &
// Example of AST:
// `-ImplicitCastExpr 0x2d0fe38 <col:9, col:10> 'int *' <BitCast>
//   `-UnaryOperator 0x2d0fe18 <col:9, col:10> 'int (*)[5]' prefix '&'
//     `-DeclRefExpr 0x2d0fdc0 <col:10> 'int [5]' lvalue Var 0x2d0fb20 'arr' 'int [5]'
func transpileUnaryOperatorAmpersant(n *ast.UnaryOperator, p *program.Program) (
	expr goast.Expr, eType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileUnaryOperatorAmpersant : err = %v", err)
		}
	}()

	expr, eType, preStmts, postStmts, err = transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return
	}
	if expr == nil {
		err = fmt.Errorf("Expr is nil")
		return
	}

	if types.IsFunction(eType) {
		return
	}

	if types.IsLastArray(eType) {
		// In : eType = 'int [5]'
		// Out: eType = 'int *'
		f := strings.Index(eType, "[")
		e := strings.Index(eType, "]")
		if e == len(eType)-1 {
			eType = eType[:f] + "*"
		} else {
			eType = eType[:f] + "*" + eType[e+1:]
		}
		expr = &goast.UnaryExpr{
			X: &goast.IndexExpr{
				X:     expr,
				Index: util.NewIntLit(0),
			},
			Op: token.AND,
		}
		return
	}

	// In : eType = 'int'
	// Out: eType = 'int *'
	// FIXME: This will need to use a real slice to reference the original
	// value.
	_, err = types.ResolveType(p, eType)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
		return
	}

	p.AddImport("unsafe")
	expr = &goast.UnaryExpr{
		X:  expr,
		Op: token.AND,
	}

	// We now have a pointer to the original type.
	eType += " *"
	return
}

// transpilePointerArith - transpile pointer arithmetic
// Example of using:
// *(t + 1) = ...
func transpilePointerArith(n *ast.UnaryOperator, p *program.Program) (
	expr goast.Expr, eType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	// pointer - expression with name of array pointer
	var pointer interface{}

	// locationPointer
	var locPointer ast.Node
	var locPosition int

	// counter - count of amount of changes in AST tree
	var counter int

	var parents []ast.Node
	var found bool

	var f func(ast.Node)
	f = func(n ast.Node) {
		for i := range n.Children() {
			switch v := n.Children()[i].(type) {
			case *ast.ArraySubscriptExpr,
				*ast.UnaryOperator,
				*ast.DeclRefExpr:
				counter++
				if counter > 1 {
					err = fmt.Errorf("Not acceptable : change counter is more then 1. found = %T,%T", pointer, v)
					return
				}
				// found pointer
				pointer = v
				// Replace pointer to zero
				var zero ast.IntegerLiteral
				zero.Type = "int"
				zero.Value = "0"
				locPointer = n
				locPosition = i
				n.Children()[i] = &zero
				found = true
				return

			case *ast.CStyleCastExpr:
				if v.Type == "int" {
					continue
				}
				counter++
				if counter > 1 {
					err = fmt.Errorf("Not acceptable : change counter is more then 1. found = %T,%T", pointer, v)
					return
				}
				// found pointer
				pointer = v
				// Replace pointer to zero
				var zero ast.IntegerLiteral
				zero.Type = "int"
				zero.Value = "0"
				locPointer = n
				locPosition = i
				n.Children()[i] = &zero
				found = true
				return

			case *ast.MemberExpr:
				// check - if member of union
				a := n.Children()[i]
				var isUnion bool
				for {
					if a == nil {
						break
					}
					if len(a.Children()) == 0 {
						break
					}
					switch vv := a.Children()[0].(type) {
					case *ast.MemberExpr, *ast.DeclRefExpr:
						var typeVV string
						if vvv, ok := vv.(*ast.MemberExpr); ok {
							typeVV = vvv.Type
						}
						if vvv, ok := vv.(*ast.DeclRefExpr); ok {
							typeVV = vvv.Type
						}
						typeVV = types.GetBaseType(typeVV)

						if _, ok := p.Structs[typeVV]; ok {
							isUnion = true
						}
						if _, ok := p.Structs["struct "+typeVV]; ok {
							isUnion = true
						}
						if strings.HasPrefix(typeVV, "union ") || strings.HasPrefix(typeVV, "struct ") {
							isUnion = true
						}
						if isUnion {
							break
						}
						a = vv
						continue
					case *ast.ImplicitCastExpr, *ast.CStyleCastExpr:
						a = vv
						continue
					}
					break
				}
				if isUnion {
					counter++
					if counter > 1 {
						err = fmt.Errorf("Not acceptable : change counter is more then 1. found = %v,%v", pointer, v)
						return
					}
					// found pointer
					pointer = v
					// Replace pointer to zero
					var zero ast.IntegerLiteral
					zero.Type = "int"
					zero.Value = "0"
					locPointer = n
					locPosition = i
					n.Children()[i] = &zero
					found = true
					return
				}
				// member of struct
				f(v)

			case *ast.CallExpr:
				if v.Type == "int" {
					continue
				}
				counter++
				if counter > 1 {
					err = fmt.Errorf("Not acceptable : change counter is more then 1. found = %T,%T", pointer, v)
					return
				}
				// found pointer
				pointer = v
				// Replace pointer to zero
				var zero ast.IntegerLiteral
				zero.Type = "int"
				zero.Value = "0"
				locPointer = n
				locPosition = i
				n.Children()[i] = &zero
				found = true
				return

			default:
				if found {
					break
				}
				if len(v.Children()) > 0 {
					if found {
						break
					}
					parents = append(parents, v)
					deep := true
					if vv, ok := v.(*ast.ImplicitCastExpr); ok && types.IsCInteger(p, vv.Type) {
						deep = false
					}
					if vv, ok := v.(*ast.CStyleCastExpr); ok && types.IsCInteger(p, vv.Type) {
						deep = false
					}
					if deep {
						f(v)
					}
					if !found {
						parents = parents[:len(parents)-1]
					}
				}
			}
		}
	}
	f(n)

	if err != nil {
		return
	}

	if pointer == nil {
		err = fmt.Errorf("pointer is nil")
		return
	}

	defer func() {
		if pointer != nil && locPointer != nil {
			locPointer.Children()[locPosition] = pointer.(ast.Node)
		}
	}()

	var typesParentBefore []string
	for i := range parents {
		switch v := parents[i].(type) {
		case *ast.ParenExpr:
			typesParentBefore = append(typesParentBefore, v.Type)
			v.Type = "int"
		case *ast.BinaryOperator:
			typesParentBefore = append(typesParentBefore, v.Type)
			v.Type = "int"
		case *ast.ImplicitCastExpr:
			typesParentBefore = append(typesParentBefore, v.Type)
			v.Type = "int"
		case *ast.CStyleCastExpr:
			typesParentBefore = append(typesParentBefore, v.Type)
			v.Type = "int"
		case *ast.VAArgExpr:
			typesParentBefore = append(typesParentBefore, v.Type)
			v.Type = "int"
		case *ast.MemberExpr:
			typesParentBefore = append(typesParentBefore, v.Type)
			v.Type = "int"
		default:
			panic(fmt.Errorf("Not support parent type %T in pointer searching", v))
		}
	}
	defer func() {
		for i := range parents {
			switch v := parents[i].(type) {
			case *ast.ParenExpr:
				v.Type = typesParentBefore[i]
			case *ast.BinaryOperator:
				v.Type = typesParentBefore[i]
			case *ast.ImplicitCastExpr:
				v.Type = typesParentBefore[i]
			case *ast.CStyleCastExpr:
				v.Type = typesParentBefore[i]
			case *ast.VAArgExpr:
				v.Type = typesParentBefore[i]
			case *ast.MemberExpr:
				v.Type = typesParentBefore[i]
			default:
				panic(fmt.Errorf("Not support parent type %T in pointer searching", v))
			}
		}
	}()

	e, eType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
	eType = n.Type

	switch v := pointer.(type) {
	case *ast.MemberExpr:
		arr, arrType, newPre, newPost, err2 := transpileToExpr(v, p, false)
		if err2 != nil {
			return
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		isConst, indexInt := util.EvaluateConstExpr(e)
		if isConst && indexInt == 0 {
			// nop
		} else if isConst && indexInt < 0 {
			indexInt = -indexInt
			arr, _, newPre, newPost, err =
				pointerArithmetic(p, arr, arrType, util.NewIntLit(int(indexInt)), "int", token.SUB)
		} else {
			arr, _, newPre, newPost, err =
				pointerArithmetic(p, arr, arrType, e, eType, token.ADD)
		}
		return &goast.StarExpr{
			X: arr,
		}, eType, preStmts, postStmts, err

	case *ast.DeclRefExpr:
		var ident goast.Expr
		ident = util.NewIdent(v.Name)
		isConst, indexInt := util.EvaluateConstExpr(e)
		if isConst && indexInt == 0 {
			if strings.HasSuffix(v.Type, "]") {
				return &goast.IndexExpr{
					X:     ident,
					Index: util.NewIntLit(0),
				}, eType, preStmts, postStmts, err
			}
		} else if isConst && indexInt < 0 {
			indexInt = -indexInt
			ident, _, newPre, newPost, err =
				pointerArithmetic(p, ident, n.Type+" *", util.NewIntLit(int(indexInt)), "int", token.SUB)
		} else {
			ident, _, newPre, newPost, err =
				pointerArithmetic(p, ident, n.Type+" *", e, eType, token.ADD)
		}

		return &goast.StarExpr{
			X: ident,
		}, eType, preStmts, postStmts, err

	case *ast.ArraySubscriptExpr, *ast.CallExpr, *ast.CStyleCastExpr:
		arr, arrType, newPre, newPost, err2 := transpileToExpr(v.(ast.Node), p, false)
		if err2 != nil {
			return
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		isConst, indexInt := util.EvaluateConstExpr(e)
		if isConst && indexInt == 0 {
			// nop
		} else if isConst && indexInt < 0 {
			indexInt = -indexInt
			arr, _, newPre, newPost, err =
				pointerArithmetic(p, arr, arrType, util.NewIntLit(int(indexInt)), "int", token.SUB)
		} else {
			arr, _, newPre, newPost, err =
				pointerArithmetic(p, arr, arrType, e, eType, token.ADD)
		}
		return &goast.StarExpr{
			X: arr,
		}, eType, preStmts, postStmts, err

	case *ast.UnaryOperator:
		arr, arrType, newPre, newPost, err2 := atomicOperation(v, p)
		if err2 != nil {
			return
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		if memberName, ok := getMemberName(n.Children()[0]); ok {
			return &goast.StarExpr{
				X: &goast.SelectorExpr{
					X:   arr,
					Sel: util.NewIdent(memberName),
				},
			}, eType, preStmts, postStmts, err
		}
		isConst, indexInt := util.EvaluateConstExpr(e)
		if isConst && indexInt == 0 {
			// nop
		} else if isConst && indexInt < 0 {
			indexInt = -indexInt
			arr, _, newPre, newPost, err =
				pointerArithmetic(p, arr, arrType, util.NewIntLit(int(indexInt)), "int", token.SUB)
		} else {
			arr, _, newPre, newPost, err =
				pointerArithmetic(p, arr, arrType, e, eType, token.ADD)
		}
		return &goast.StarExpr{
			X: arr,
		}, eType, preStmts, postStmts, err
	}
	return nil, "", nil, nil, fmt.Errorf("Cannot found : %#v", pointer)
}

func transpileUnaryOperator(n *ast.UnaryOperator, p *program.Program, exprIsStmt bool) (
	_ goast.Expr, theType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile UnaryOperator: err = %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	operator := getTokenForOperator(n.Operator)

	switch operator {
	case token.MUL: // *
		// Prefix "*" is not a multiplication.
		// Prefix "*" used for pointer arithmetic
		// Example of using:
		// *(t + 1) = ...
		return transpilePointerArith(n, p)
	case token.INC, token.DEC: // ++, --
		return transpileUnaryOperatorInc(n, p, operator, exprIsStmt)
	case token.NOT: // !
		return transpileUnaryOperatorNot(n, p)
	case token.AND: // &
		return transpileUnaryOperatorAmpersant(n, p)
	}

	// Otherwise handle like a unary operator.
	e, eType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, exprIsStmt)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	return &goast.UnaryExpr{
		Op: operator,
		X:  e,
	}, eType, preStmts, postStmts, nil

}

func transpileUnaryExprOrTypeTraitExpr(n *ast.UnaryExprOrTypeTraitExpr, p *program.Program) (
	*goast.BasicLit, string, []goast.Stmt, []goast.Stmt, error) {
	t := n.Type2

	// It will have children if the sizeof() is referencing a variable.
	// Fortunately clang already has the type in the AST for us.
	if len(n.Children()) > 0 {
		var realFirstChild interface{}
		t = ""

		switch c := n.Children()[0].(type) {
		case *ast.ParenExpr:
			realFirstChild = c.Children()[0]
		case *ast.DeclRefExpr:
			t = c.Type
		default:
			panic(fmt.Sprintf("cannot find first child from: %#v", n.Children()[0]))
		}

		if t == "" {
			switch ty := realFirstChild.(type) {
			case *ast.DeclRefExpr:
				t = ty.Type2

			case *ast.ArraySubscriptExpr:
				t = ty.Type

			case *ast.MemberExpr:
				t = ty.Type

			case *ast.UnaryOperator:
				t = ty.Type

			case *ast.ParenExpr:
				t = ty.Type

			case *ast.CallExpr:
				t = ty.Type

			default:
				panic(fmt.Sprintf("cannot do unary on: %#v", ty))
			}
		}
	}

	sizeInBytes, err := types.SizeOf(p, t)
	p.AddMessage(p.GenerateWarningMessage(err, n))

	return util.NewIntLit(sizeInBytes), n.Type1, nil, nil, nil
}

func transpileStmtExpr(n *ast.StmtExpr, p *program.Program) (
	*goast.CallExpr, string, []goast.Stmt, []goast.Stmt, error) {
	returnType, err := types.ResolveType(p, n.Type)
	if err != nil {
		return nil, "", nil, nil, err
	}

	body, pre, post, err := transpileCompoundStmt(n.Children()[0].(*ast.CompoundStmt), p)
	if err != nil {
		return nil, "", pre, post, err
	}

	// The body of the StmtExpr is always a CompoundStmt. However, the last
	// statement needs to be transformed into an explicit return statement.
	lastStmt := body.List[len(body.List)-1]
	body.List[len(body.List)-1] = &goast.ReturnStmt{
		Results: []goast.Expr{lastStmt.(*goast.ExprStmt).X},
	}

	return util.NewFuncClosure(returnType, body.List...), n.Type, pre, post, nil
}
