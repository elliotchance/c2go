// This file contains functions for transpiling unary operator expressions.

package transpiler

import (
	"fmt"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"

	goast "go/ast"
	"go/token"
)

func transpileUnaryOperator(n *ast.UnaryOperator, p *program.Program) (
	goast.Expr, string, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}
	operator := getTokenForOperator(n.Operator)

	// Unfortunately we cannot use the Go increment operators because we are not
	// providing any position information for tokens. This means that the ++/--
	// would be placed before the expression and would be invalid in Go.
	//
	// Until it can be properly fixed (can we trick Go into to placing it after
	// the expression with a magic position?) we will have to return a
	// BinaryExpr with the same functionality.
	if operator == token.INC || operator == token.DEC {
		binaryOperator := "+="
		if operator == token.DEC {
			binaryOperator = "-="
		}

		return transpileBinaryOperator(&ast.BinaryOperator{
			Type:     n.Type,
			Operator: binaryOperator,
			Children: []ast.Node{
				n.Children[0], &ast.IntegerLiteral{
					Type:     "int",
					Value:    "1",
					Children: []ast.Node{},
				},
			},
		}, p)
	}

	// Otherwise handle like a unary operator.
	e, eType, newPre, newPost, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	if operator == token.NOT {
		if eType == "bool" || eType == "_Bool" {
			return &goast.UnaryExpr{
				X:  e,
				Op: operator,
			}, "bool", preStmts, postStmts, nil
		}

		t, err := types.ResolveType(p, eType)
		ast.IsWarning(err, n)

		if t == "string" {
			return &goast.BinaryExpr{
				X:  e,
				Op: token.EQL,
				Y: &goast.BasicLit{
					Kind:  token.STRING,
					Value: `""`,
				},
			}, "bool", preStmts, postStmts, nil
		}

		p.AddImport("github.com/elliotchance/c2go/noarch")

		functionName := fmt.Sprintf("noarch.Not%s", util.Ucfirst(t))

		return &goast.CallExpr{
			Fun:  goast.NewIdent(functionName),
			Args: []goast.Expr{e},
		}, eType, preStmts, postStmts, nil
	}

	if operator == token.MUL {
		if eType == "const char *" {
			return &goast.IndexExpr{
				X: e,
				Index: &goast.BasicLit{
					Kind:  token.INT,
					Value: "0",
				},
			}, "char", preStmts, postStmts, nil
		}

		t, err := types.GetDereferenceType(eType)
		if err != nil {
			return nil, "", preStmts, postStmts, err
		}

		return &goast.StarExpr{
			X: e,
		}, t, preStmts, postStmts, nil
	}

	if operator == token.AND {
		eType += " *"
	}

	return &goast.UnaryExpr{
		X:  e,
		Op: operator,
	}, eType, preStmts, postStmts, nil
}

func transpileUnaryExprOrTypeTraitExpr(n *ast.UnaryExprOrTypeTraitExpr, p *program.Program) (
	*goast.BasicLit, string, []goast.Stmt, []goast.Stmt, error) {
	t := n.Type2

	// It will have children if the sizeof() is referencing a variable.
	// Fortunately clang already has the type in the AST for us.
	if len(n.Children) > 0 {
		switch ty := n.Children[0].(*ast.ParenExpr).Children[0].(type) {
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

		default:
			panic(fmt.Sprintf("cannot do unary on: %#v", ty))
		}
	}

	ty, err := types.ResolveType(p, n.Type1)
	ast.IsWarning(err, n)

	sizeInBytes, err := types.SizeOf(p, t)
	ast.IsWarning(err, n)

	return util.NewIntLit(sizeInBytes), ty, nil, nil, nil
}
