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

func transpileUnaryOperatorInc(n *ast.UnaryOperator, p *program.Program,
	operator token.Token) (goast.Expr, string, []goast.Stmt, []goast.Stmt, error) {
	// Unfortunately we cannot use the Go increment operators because we are not
	// providing any position information for tokens. This means that the ++/--
	// would be placed before the expression and would be invalid in Go.
	//
	// Until it can be properly fixed (can we trick Go into to placing it after
	// the expression with a magic position?) we will have to return a
	// BinaryExpr with the same functionality.

	// Construct code for assigning value to an union field
	memberExpr, ok := n.Children()[0].(*ast.MemberExpr)
	if ok {
		ref := memberExpr.GetDeclRefExpr()
		if ref != nil {
			binaryOperator := token.ADD
			if operator == token.DEC {
				binaryOperator = token.SUB
			}

			// TODO: Is this duplicate code in operators.go?
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
				argOp := binaryOperator
				argRHS := util.NewIntLit(1)
				argValue := util.NewBinaryExpr(argLHS, argOp, argRHS, "interface{}", false)

				// Make Go expression
				resExpr := util.NewCallExpr(setterName, argValue)

				return resExpr, n.Type, nil, nil, nil
			}
		}
	}

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
	}, p, false)
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

	if t == "[]byte" {
		return util.NewUnaryExpr(
			token.NOT, util.NewCallExpr("noarch.CStringIsNull", e),
		), "bool", preStmts, postStmts, nil
	}

	p.AddImport("github.com/elliotchance/c2go/noarch")

	functionName := fmt.Sprintf("noarch.Not%s",
		util.GetExportedName(t))

	return util.NewCallExpr(functionName, e),
		eType, preStmts, postStmts, nil
}

// transpilePointerArith - transpile pointer aripthmetic
// Example of using:
// *(t + 1) = ...
func transpilePointerArith(n *ast.UnaryOperator, p *program.Program) (
	expr goast.Expr, eType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	// pointer - expression with name of array pointer
	var pointer interface{}

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
					err = fmt.Errorf("Not acceptable : change counter is more then 1. found = %v,%v", pointer, v)
					return
				}
				// found pointer
				pointer = v
				// Replace pointer to zero
				var zero ast.IntegerLiteral
				zero.Type = "int"
				zero.Value = "0"
				n.Children()[i] = &zero
				found = true
				return

			case *ast.CStyleCastExpr:
				if v.Type == "int" {
					continue
				}
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
				n.Children()[i] = &zero
				found = true
				return

			case *ast.CallExpr:
				if v.Type == "int" {
					continue
				}
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
					var deep bool = true
					if vv, ok := v.(*ast.ImplicitCastExpr); ok && vv.Type == "int" {
						deep = false
					}
					if vv, ok := v.(*ast.CStyleCastExpr); ok && vv.Type == "int" {
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
		err = fmt.Errorf("pointer of array is nil")
		return
	}
	for i := range parents {
		switch v := parents[i].(type) {
		case *ast.ParenExpr:
			v.Type = "int"
		case *ast.BinaryOperator:
			v.Type = "int"
		case *ast.ImplicitCastExpr:
			v.Type = "int"
		case *ast.CStyleCastExpr:
			v.Type = "int"
		case *ast.VAArgExpr:
			v.Type = "int"
		case *ast.MemberExpr:
			v.Type = "int"
		default:
			panic(fmt.Errorf("Not support parent type %T in pointer seaching", v))
		}
	}

	e, eType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
	eType = n.Type

	switch v := pointer.(type) {
	case *ast.DeclRefExpr:
		return &goast.IndexExpr{
			X:     util.NewIdent(v.Name),
			Index: e,
		}, eType, preStmts, postStmts, err
	case *ast.ArraySubscriptExpr:
		arr, _, newPre, newPost, err2 := transpileToExpr(v, p, false)
		if err2 != nil {
			return
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		return &goast.IndexExpr{
			X: &goast.ParenExpr{
				Lparen: 1,
				X:      arr,
			},
			Index: e,
		}, eType, preStmts, postStmts, err
	case *ast.CallExpr:
		arr, _, newPre, newPost, err2 := transpileToExpr(v, p, false)
		if err2 != nil {
			return
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		return &goast.IndexExpr{
			X: &goast.ParenExpr{
				Lparen: 1,
				X:      arr,
			},
			Index: e,
		}, eType, preStmts, postStmts, err
	case *ast.CStyleCastExpr:
		arr, _, newPre, newPost, err2 := transpileToExpr(v, p, false)
		if err2 != nil {
			return
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		return &goast.IndexExpr{
			X: &goast.ParenExpr{
				Lparen: 1,
				X:      arr,
			},
			Index: e,
		}, eType, preStmts, postStmts, err
	case *ast.UnaryOperator:
		arr, _, newPre, newPost, err2 := transpileToExpr(v, p, false)
		if err2 != nil {
			return
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		return &goast.IndexExpr{
			X: &goast.ParenExpr{
				Lparen: 1,
				X:      arr,
			},
			Index: e,
		}, eType, preStmts, postStmts, err
	}
	return nil, "", nil, nil, fmt.Errorf("Cannot found : %#v", pointer)
}
func transpileUnaryOperator(n *ast.UnaryOperator, p *program.Program) (
	_ goast.Expr, theType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile UnaryOperator: err = %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	operator := getTokenForOperator(n.Operator)

	// Prefix "*" is not a multiplication.
	// Prefix "*" used for pointer ariphmetic
	// Example of using:
	// *(t + 1) = ...
	if n.IsPrefix && n.IsLvalue && n.Operator == "*" {
		expr, eType, preStmts, postStmts, err := transpilePointerArith(n, p)
		if err != nil {
			return nil, "", nil, nil, err
		}
		return expr, eType, preStmts, postStmts, nil
	}

	switch operator {
	case token.INC, token.DEC:
		return transpileUnaryOperatorInc(n, p, operator)
	case token.NOT:
		return transpileUnaryOperatorNot(n, p)
	}

	// Otherwise handle like a unary operator.
	e, eType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	if operator == token.AND {
		// FIXME: This will need to use a real slice to reference the original
		// value.
		resolvedType, err := types.ResolveType(p, eType)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}

		p.AddImport("unsafe")
		e = util.CreateSliceFromReference(resolvedType, e)

		// We now have a pointer to the original type.
		eType += " *"

		return e, eType, preStmts, postStmts, nil
	}

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
