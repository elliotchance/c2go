package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
)

func TranspileAST(fileName string, p *program.Program, root ast.Node) error {
	// Start by parsing an empty file.
	p.FileSet = token.NewFileSet()
	f, err := parser.ParseFile(p.FileSet, fileName, "package main", 0)
	p.File = f

	if err != nil {
		return err
	}

	// Now begin building the Go AST.
	err = transpileToNode(root, p)
	return err
}

func transpileToExpr(node ast.Node, p *program.Program) (goast.Expr, string, error) {
	if node == nil {
		return nil, "", nil
	}

	switch n := node.(type) {
	case *ast.StringLiteral:
		return transpileStringLiteral(n), "", nil

	case *ast.FloatingLiteral:
		return transpileFloatingLiteral(n), "", nil

	case *ast.PredefinedExpr:
		return transpilePredefinedExpr(n, p)

	case *ast.ConditionalOperator:
		return transpileConditionalOperator(n, p)

	case *ast.ArraySubscriptExpr:
		return transpileArraySubscriptExpr(n, p)

	case *ast.BinaryOperator:
		return transpileBinaryOperator(n, p)

	case *ast.UnaryOperator:
		return transpileUnaryOperator(n, p)

	case *ast.MemberExpr:
		return transpileMemberExpr(n, p)

	case *ast.ImplicitCastExpr:
		return transpileToExpr(n.Children[0], p)

	case *ast.DeclRefExpr:
		return transpileDeclRefExpr(n, p)

	case *ast.IntegerLiteral:
		return transpileIntegerLiteral(n), "", nil

	case *ast.ParenExpr:
		return transpileParenExpr(n, p)

	case *ast.CStyleCastExpr:
		return transpileToExpr(n.Children[0], p)

	case *ast.CharacterLiteral:
		return transpileCharacterLiteral(n), "", nil

	case *ast.CallExpr:
		return transpileCallExpr(n, p)
	}

	panic(fmt.Sprintf("cannot transpile to expr: %#v", node))
}

func transpileCompoundStmt(n *ast.CompoundStmt, p *program.Program) (*goast.BlockStmt, error) {
	stmts := []goast.Stmt{}

	for _, x := range n.Children {
		result, err := transpileToStmt(x, p)
		if err != nil {
			return nil, err
		}

		if result != nil {
			stmts = append(stmts, result)
		}
	}

	return &goast.BlockStmt{
		Lbrace: token.NoPos,
		List:   stmts,
		Rbrace: token.NoPos,
	}, nil
}

func transpileToStmt(node ast.Node, p *program.Program) (goast.Stmt, error) {
	if node == nil {
		return nil, nil
	}

	switch n := node.(type) {
	case *ast.DefaultStmt:
		return transpileDefaultStmt(n, p)

	case *ast.CaseStmt:
		return transpileCaseStmt(n, p)

	case *ast.SwitchStmt:
		return transpileSwitchStmt(n, p)

	case *ast.BreakStmt:
		return &goast.BranchStmt{
			Tok: token.BREAK,
		}, nil

	case *ast.WhileStmt:
		return transpileWhileStmt(n, p)

	case *ast.IfStmt:
		return transpileIfStmt(n, p)

	case *ast.ForStmt:
		return transpileForStmt(n, p)

	case *ast.DeclStmt:
		return transpileDeclStmt(n, p)

	case *ast.ReturnStmt:
		return transpileReturnStmt(n, p)

	case *ast.CompoundStmt:
		return transpileCompoundStmt(n, p)
	}

	e, _, err := transpileToExpr(node, p)
	if err != nil {
		return nil, err
	}

	return &goast.ExprStmt{
		X: e,
	}, nil
}

func transpileToBlockStmt(node ast.Node, p *program.Program) (*goast.BlockStmt, error) {
	e, err := transpileToStmt(node, p)
	if err != nil {
		return nil, err
	}

	if block, ok := e.(*goast.BlockStmt); ok {
		return block, nil
	}

	return &goast.BlockStmt{
		List: []goast.Stmt{e},
	}, nil
}

func transpileToNode(node ast.Node, p *program.Program) error {
	switch n := node.(type) {
	case *ast.TranslationUnitDecl:
		for _, c := range n.Children {
			transpileToNode(c, p)
		}

	case *ast.FunctionDecl:
		err := transpileFunctionDecl(n, p)
		if err != nil {
			return err
		}

	case *ast.TypedefDecl, *ast.RecordDecl, *ast.VarDecl:
		// do nothing

	default:
		panic(fmt.Sprintf("cannot transpile to node: %#v", node))
	}

	return nil
}

func transpileStmts(nodes []ast.Node, p *program.Program) ([]goast.Stmt, error) {
	stmts := []goast.Stmt{}

	for _, s := range nodes {
		if s != nil {
			a, err := transpileToStmt(s, p)
			if err != nil {
				return nil, err
			}

			stmts = append(stmts, a)
		}
	}

	return stmts, nil
}
