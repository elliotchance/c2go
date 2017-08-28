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

	if eType == "bool" || eType == "_Bool" {
		return &goast.UnaryExpr{
			X:  e,
			Op: token.NOT,
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

// Dereferencing.
func transpileUnaryOperatorMul(n *ast.UnaryOperator, p *program.Program) (
	goast.Expr, string, []goast.Stmt, []goast.Stmt, error) {
	e, eType, preStmts, postStmts, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	if eType == "const char *" {
		return &goast.IndexExpr{
			X:     e,
			Index: util.NewIntLit(0),
		}, "char", preStmts, postStmts, nil
	}

	t, err := types.GetDereferenceType(eType)
	if err != nil {
		return nil, "", preStmts, postStmts, err
	}

	// C is more relaxed with this syntax. In Go we convert all of the
	// pointers to slices, so we have to be careful when dereference a slice
	// that it actually takes the first element instead.
	resolvedType, err := types.ResolveType(p, eType)
	if strings.HasPrefix(resolvedType, "[]") {
		return &goast.IndexExpr{
			X:     e,
			Index: util.NewIntLit(0),
		}, t, preStmts, postStmts, nil
	}

	return &goast.StarExpr{
		X: e,
	}, t, preStmts, postStmts, nil
}

func transpileUnaryOperator(n *ast.UnaryOperator, p *program.Program) (
	goast.Expr, string, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}
	operator := getTokenForOperator(n.Operator)

	switch operator {
	case token.INC, token.DEC:
		return transpileUnaryOperatorInc(n, p, operator)
	case token.NOT:
		return transpileUnaryOperatorNot(n, p)
	case token.MUL:
		return transpileUnaryOperatorMul(n, p)
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
