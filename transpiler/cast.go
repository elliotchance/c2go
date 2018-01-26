package transpiler

import (
	"fmt"
	goast "go/ast"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
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
		expr = util.NewIdent("nil")
		exprType = types.NullPointer
		return
	}
	if strings.Contains(n.Type, "enum") {
		if d, ok := n.Children()[0].(*ast.DeclRefExpr); ok {
			expr, exprType, err = util.NewIdent(d.Name), n.Type, nil
			return
		}
	}
	expr, exprType, preStmts, postStmts, err = transpileToExpr(n.Children()[0], p, exprIsStmt)
	if err != nil {
		return nil, "", nil, nil, err
	}
	if exprType == types.NullPointer {
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

	if n.Kind == ast.CStyleCastExprNullToPointer {
		expr = util.NewIdent("nil")
		exprType = types.NullPointer
		return
	}
	expr, exprType, preStmts, postStmts, err = transpileToExpr(n.Children()[0], p, exprIsStmt)
	if err != nil {
		return nil, "", nil, nil, err
	}

	if exprType == types.NullPointer {
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
