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

	// Add the imports after everything else so we can ensure that they are all
	// placed at the top.
	for _, quotedImportPath := range p.Imports() {
		importSpec := &goast.ImportSpec{
			Path: &goast.BasicLit{
				Kind:  token.IMPORT,
				Value: quotedImportPath,
			},
		}
		importDecl := &goast.GenDecl{
			Tok: token.IMPORT,
		}

		importDecl.Specs = append(importDecl.Specs, importSpec)
		p.File.Decls = append([]goast.Decl{importDecl}, p.File.Decls...)
	}

	return err
}

func transpileToExpr(node ast.Node, p *program.Program) (goast.Expr, string, error) {
	if node == nil {
		panic(node)
		return nil, "unknown1", nil
	}

	switch n := node.(type) {
	// TODO: It would be much easier if all of these returned three arguments.
	case *ast.StringLiteral:
		return transpileStringLiteral(n), "const char *", nil

	case *ast.FloatingLiteral:
		return transpileFloatingLiteral(n), "double", nil

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
		return transpileIntegerLiteral(n), "int", nil

	case *ast.ParenExpr:
		return transpileParenExpr(n, p)

	case *ast.CStyleCastExpr:
		return transpileToExpr(n.Children[0], p)

	case *ast.CharacterLiteral:
		return transpileCharacterLiteral(n), "char", nil

	case *ast.CallExpr:
		return transpileCallExpr(n, p)
	}

	panic(fmt.Sprintf("cannot transpile to expr: %#v", node))
}

func transpileToStmts(node ast.Node, p *program.Program) ([]goast.Stmt, error) {
	if node == nil {
		return nil, nil
	}

	switch n := node.(type) {
	case *ast.DeclStmt:
		return transpileDeclStmt(n, p)
	}

	stmt, err := transpileToStmt(node, p)
	return []goast.Stmt{stmt}, err
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

	case *ast.DoStmt:
		return transpileDoStmt(n, p)

	case *ast.IfStmt:
		return transpileIfStmt(n, p)

	case *ast.ForStmt:
		return transpileForStmt(n, p)

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

	case *ast.TypedefDecl:
		return transpileTypedefDecl(p, n)

	case *ast.RecordDecl:
		return transpileRecordDecl(p, n)

	case *ast.VarDecl:
		transpileVarDecl(p, n)
		return nil

	case *ast.EnumDecl:
		transpileEnumDecl(p, n)
		return nil

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
