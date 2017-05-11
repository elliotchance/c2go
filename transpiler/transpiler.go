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

	// Now we need to build the __init() function. This sets up certain state
	// and variables that the runtime expects to be ready.
	p.File.Decls = append(p.File.Decls, &goast.FuncDecl{
		Name: goast.NewIdent("__init"),
		Type: &goast.FuncType{
			Params: &goast.FieldList{
				List: []*goast.Field{},
			},
			Results: nil,
		},
		Body: &goast.BlockStmt{
			List: p.StartupStatements(),
		},
	})

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

func transpileToExpr(node ast.Node, p *program.Program) (
	expr goast.Expr,
	exprType string,
	preStmts []goast.Stmt,
	postStmts []goast.Stmt,
	err error) {
	if node == nil {
		panic(node)
	}

	switch n := node.(type) {
	case *ast.StringLiteral:
		expr = transpileStringLiteral(n)
		exprType = "const char *"

	case *ast.FloatingLiteral:
		expr = transpileFloatingLiteral(n)
		exprType = "double"

	case *ast.PredefinedExpr:
		expr, exprType, err = transpilePredefinedExpr(n, p)

	case *ast.ConditionalOperator:
		expr, exprType, preStmts, postStmts, err = transpileConditionalOperator(n, p)

	case *ast.ArraySubscriptExpr:
		expr, exprType, preStmts, postStmts, err = transpileArraySubscriptExpr(n, p)

	case *ast.BinaryOperator:
		expr, exprType, preStmts, postStmts, err = transpileBinaryOperator(n, p)

	case *ast.UnaryOperator:
		expr, exprType, preStmts, postStmts, err = transpileUnaryOperator(n, p)

	case *ast.MemberExpr:
		expr, exprType, preStmts, postStmts, err = transpileMemberExpr(n, p)

	case *ast.ImplicitCastExpr:
		expr, exprType, preStmts, postStmts, err = transpileToExpr(n.Children[0], p)

	case *ast.DeclRefExpr:
		expr, exprType, err = transpileDeclRefExpr(n, p)

	case *ast.IntegerLiteral:
		expr, exprType, err = transpileIntegerLiteral(n), "int", nil

	case *ast.ParenExpr:
		expr, exprType, preStmts, postStmts, err = transpileParenExpr(n, p)

	case *ast.CStyleCastExpr:
		expr, exprType, preStmts, postStmts, err = transpileToExpr(n.Children[0], p)

	case *ast.CharacterLiteral:
		expr, exprType, err = transpileCharacterLiteral(n), "char", nil

	case *ast.CallExpr:
		expr, exprType, preStmts, postStmts, err = transpileCallExpr(n, p)

	default:
		panic(fmt.Sprintf("cannot transpile to expr: %#v", node))
	}

	// Real return is through named arguments.
	return
}

func transpileToStmts(node ast.Node, p *program.Program) ([]goast.Stmt, error) {
	if node == nil {
		return nil, nil
	}

	switch n := node.(type) {
	case *ast.DeclStmt:
		stmts, preStmts, postStmts, err := transpileDeclStmt(n, p)
		stmts = append(preStmts, stmts...)
		stmts = append(stmts, postStmts...)
		return stmts, err
	}

	stmt, preStmts, postStmts, err := transpileToStmt(node, p)
	stmts := append(preStmts, stmt)
	stmts = append(stmts, postStmts...)
	return stmts, err
}

func transpileToStmt(node ast.Node, p *program.Program) (
	stmt goast.Stmt, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	if node == nil {
		return
	}

	var expr goast.Expr

	switch n := node.(type) {
	case *ast.DefaultStmt:
		stmt, err = transpileDefaultStmt(n, p)
		return

	case *ast.CaseStmt:
		stmt, preStmts, postStmts, err = transpileCaseStmt(n, p)
		return

	case *ast.SwitchStmt:
		stmt, preStmts, postStmts, err = transpileSwitchStmt(n, p)
		return

	case *ast.BreakStmt:
		stmt = &goast.BranchStmt{
			Tok: token.BREAK,
		}
		return

	case *ast.WhileStmt:
		return transpileWhileStmt(n, p)

	case *ast.DoStmt:
		return transpileDoStmt(n, p)

	case *ast.ContinueStmt:
		stmt, err = transpileContinueStmt(n, p)
		return

	case *ast.IfStmt:
		return transpileIfStmt(n, p)

	case *ast.ForStmt:
		return transpileForStmt(n, p)

	case *ast.ReturnStmt:
		return transpileReturnStmt(n, p)

	case *ast.CompoundStmt:
		stmt, preStmts, postStmts, err = transpileCompoundStmt(n, p)
		return
	}

	// We do not care about the return type.
	expr, _, preStmts, postStmts, err = transpileToExpr(node, p)
	if err != nil {
		return
	}

	stmt = &goast.ExprStmt{
		X: expr,
	}

	return
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
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}
	stmts := []goast.Stmt{}

	for _, s := range nodes {
		if s != nil {
			a, newPre, newPost, err := transpileToStmt(s, p)
			if err != nil {
				return nil, err
			}

			preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

			stmts = append(stmts, a)
		}
	}

	return stmts, nil
}
