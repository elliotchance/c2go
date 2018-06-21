package transpiler

import (
	"fmt"
	goast "go/ast"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
	"go/token"
)

func transpileImplicitCastExpr(n *ast.ImplicitCastExpr, p *program.Program, exprIsStmt bool) (
	expr goast.Expr,
	exprType string,
	preStmts []goast.Stmt,
	postStmts []goast.Stmt,
	err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileImplicitCastExpr. err = %v", err)
		}
	}()

	if n.Kind == ast.CStyleCastExprNullToPointer {
		expr = goast.NewIdent("nil")
		exprType = types.NullPointer
		return
	}
	if strings.Contains(n.Type, "enum") {
		if d, ok := n.Children()[0].(*ast.DeclRefExpr); ok {
			expr, exprType, err = util.NewIdent(d.Name), n.Type, nil
			return
		}
	}
	if isCastToUnsignedOfUnaryComplement(n, p) {
		return swapCastAndComplement(n, p, exprIsStmt)
	}
	expr, exprType, preStmts, postStmts, err = transpileToExpr(n.Children()[0], p, exprIsStmt)
	if err != nil {
		return nil, "", nil, nil, err
	}
	if exprType == types.NullPointer {
		expr = goast.NewIdent("nil")
		return
	}

	if len(n.Type) != 0 && len(n.Type2) != 0 && n.Type != n.Type2 {
		var tt string
		tt, err = types.ResolveType(p, n.Type)
		expr = &goast.CallExpr{
			Fun:    goast.NewIdent(tt),
			Lparen: 1,
			Args:   []goast.Expr{expr},
		}
		exprType = n.Type
		return
	}

	if !types.IsFunction(exprType) && !strings.ContainsAny(n.Type, "[]") {
		expr, err = types.CastExpr(p, expr, exprType, n.Type)
		if err != nil {
			return nil, "", nil, nil, err
		}
		exprType = n.Type
	}
	return
}

func isCastToUnsignedOfUnaryComplement(n *ast.ImplicitCastExpr, p *program.Program) (ret bool) {
	if !types.IsCInteger(p, n.Type) || !strings.Contains(n.Type, "unsigned ") {
		return
	}
	cn, ok := n.Children()[0].(*ast.UnaryOperator)
	if !ok || getTokenForOperator(cn.Operator) != token.XOR {
		return
	}
	return types.IsCInteger(p, cn.Type) && !strings.Contains(cn.Type, "unsigned ")
}

func swapCastAndComplement(n *ast.ImplicitCastExpr, p *program.Program, exprIsStmt bool) (
	expr goast.Expr,
	exprType string,
	preStmts []goast.Stmt,
	postStmts []goast.Stmt,
	err error) {
	uo := n.Children()[0].(*ast.UnaryOperator)
	copyUnary := &ast.UnaryOperator{}
	copyImplicit := &ast.ImplicitCastExpr{}
	*copyUnary = *uo
	*copyImplicit = *n
	unaryChildren := uo.ChildNodes
	copyUnary.ChildNodes = []ast.Node{copyImplicit}
	copyUnary.Type = copyImplicit.Type
	copyImplicit.ChildNodes = unaryChildren
	return transpileToExpr(copyUnary, p, exprIsStmt)
}

func transpileCStyleCastExpr(n *ast.CStyleCastExpr, p *program.Program, exprIsStmt bool) (
	expr goast.Expr,
	exprType string,
	preStmts []goast.Stmt,
	postStmts []goast.Stmt,
	err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileImplicitCastExpr. err = %v", err)
		}
	}()
	// Char overflow
	// example for byte(-1)
	// CStyleCastExpr 0x365f628 <col:12, col:23> 'char' <IntegralCast>
	// `-ParenExpr 0x365f608 <col:18, col:23> 'int'
	//   `-ParenExpr 0x365f5a8 <col:19, col:22> 'int'
	//     `-UnaryOperator 0x365f588 <col:20, col:21> 'int' prefix '-'
	//       `-IntegerLiteral 0x365f568 <col:21> 'int' 1
	if n.Type == "char" {
		if par, ok := n.Children()[0].(*ast.ParenExpr); ok {
			if par2, ok := par.Children()[0].(*ast.ParenExpr); ok {
				if u, ok := par2.Children()[0].(*ast.UnaryOperator); ok && u.IsPrefix {
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

	if n.Kind == ast.CStyleCastExprNullToPointer {
		expr = goast.NewIdent("nil")
		exprType = types.NullPointer
		return
	}
	expr, exprType, preStmts, postStmts, err = transpileToExpr(n.Children()[0], p, exprIsStmt)
	if err != nil {
		return nil, "", nil, nil, err
	}

	if exprType == types.NullPointer {
		expr = goast.NewIdent("nil")
		return
	}

	if len(n.Type) != 0 && len(n.Type2) != 0 && n.Type != n.Type2 {
		var tt string
		tt, err = types.ResolveType(p, n.Type)
		expr = &goast.CallExpr{
			Fun:    goast.NewIdent(tt),
			Lparen: 1,
			Args:   []goast.Expr{expr},
		}
		exprType = n.Type
		return
	}

	if n.Kind == ast.CStyleCastExprToVoid {
		exprType = types.ToVoid
		return
	}

	if !types.IsFunction(exprType) && n.Kind != ast.ImplicitCastExprArrayToPointerDecay {
		expr, err = types.CastExpr(p, expr, exprType, n.Type)
		if err != nil {
			return nil, "", nil, nil, err
		}
		exprType = n.Type
	}
	return
}
