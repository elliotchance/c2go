// Package transpiler handles the conversion between the Clang AST and the Go
// AST.
package transpiler

import (
	"errors"
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"
)

// TranspileAST iterates through the Clang AST and builds a Go AST
func TranspileAST(fileName, packageName string, p *program.Program, root ast.Node) error {
	// Start by parsing an empty file.
	p.FileSet = token.NewFileSet()
	packageSignature := fmt.Sprintf("package %v", packageName)
	f, err := parser.ParseFile(p.FileSet, fileName, packageSignature, 0)
	p.File = f

	if err != nil {
		return err
	}

	// Now begin building the Go AST.
	err = transpileToNode(root, p)

	if p.OutputAsTest {
		p.AddImport("testing")
		p.AddImport("io/ioutil")

		// TODO: There should be a cleaner way to add a function to the program.
		// This code was taken from the end of transpileFunctionDecl.
		p.File.Decls = append(p.File.Decls, &goast.FuncDecl{
			Name: util.NewIdent("TestApp"),
			Type: &goast.FuncType{
				Params: &goast.FieldList{
					List: []*goast.Field{
						&goast.Field{
							Names: []*goast.Ident{util.NewIdent("t")},
							Type:  util.NewTypeIdent("*testing.T"),
						},
					},
				},
			},
			Body: &goast.BlockStmt{
				List: []goast.Stmt{
					util.NewExprStmt(&goast.Ident{Name: "os.Chdir(\"../../..\")"}),

					// "go test" does not redirect stdin to the executable
					// running the test so we need to override them in the test
					// itself. See documentation for noarch.Stdin.
					util.NewExprStmt(&goast.Ident{Name: "ioutil.WriteFile(\"build/stdin\", []byte{'7'}, 0777)"}),
					util.NewExprStmt(
						&goast.Ident{Name: "stdin, _ := os.Open(\"build/stdin\")"},
					),
					util.NewExprStmt(util.NewBinaryExpr(
						&goast.Ident{Name: "noarch.Stdin"},
						token.ASSIGN,
						&goast.Ident{Name: "noarch.NewFile(stdin)"},
						"*noarch.File",
						true,
					)),

					util.NewExprStmt(util.NewCallExpr("main")),
				},
			},
		})
	}

	// Now we need to build the __init() function. This sets up certain state
	// and variables that the runtime expects to be ready.
	p.File.Decls = append(p.File.Decls, &goast.FuncDecl{
		Name: util.NewIdent("__init"),
		Type: util.NewFuncType(&goast.FieldList{}, ""),
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

func transpileToExpr(node ast.Node, p *program.Program, exprIsStmt bool) (
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
		expr, exprType, preStmts, postStmts, err = transpileBinaryOperator(n, p, exprIsStmt)

	case *ast.UnaryOperator:
		expr, exprType, preStmts, postStmts, err = transpileUnaryOperator(n, p)

	case *ast.MemberExpr:
		expr, exprType, preStmts, postStmts, err = transpileMemberExpr(n, p)

	case *ast.ImplicitCastExpr:
		if strings.Contains(n.Type, "enum") {
			if d, ok := n.Children()[0].(*ast.DeclRefExpr); ok {
				expr, exprType, err = util.NewIdent(d.Name), n.Type, nil
				return
			}
		}
		expr, exprType, preStmts, postStmts, err = transpileToExpr(n.Children()[0], p, exprIsStmt)

	case *ast.DeclRefExpr:
		if n.For == "EnumConstant" {
			// clang don`t show enum constant with enum type,
			// so we have to use hack for repair the type
			if v, ok := p.EnumConstantToEnum[n.Name]; ok {
				expr, exprType, err = util.NewIdent(n.Name), v, nil
				return
			}
		}
		expr, exprType, err = transpileDeclRefExpr(n, p)

	case *ast.IntegerLiteral:
		expr, exprType, err = transpileIntegerLiteral(n), "int", nil

	case *ast.ParenExpr:
		expr, exprType, preStmts, postStmts, err = transpileParenExpr(n, p)

	case *ast.CStyleCastExpr:
		expr, exprType, preStmts, postStmts, err = transpileToExpr(n.Children()[0], p, exprIsStmt)

	case *ast.CharacterLiteral:
		expr, exprType, err = transpileCharacterLiteral(n), "char", nil

	case *ast.CallExpr:
		expr, exprType, preStmts, postStmts, err = transpileCallExpr(n, p)

	case *ast.CompoundAssignOperator:
		return transpileCompoundAssignOperator(n, p, exprIsStmt)

	case *ast.UnaryExprOrTypeTraitExpr:
		return transpileUnaryExprOrTypeTraitExpr(n, p)

	case *ast.StmtExpr:
		return transpileStmtExpr(n, p)

	default:
		p.AddMessage(p.GenerateWarningMessage(errors.New("cannot transpile to expr"), node))
		expr = util.NewNil()
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

	case *ast.BinaryOperator:
		if n.Operator == "," {
			stmt, preStmts, err = transpileBinaryOperatorComma(n, p)
			return
		}

	case *ast.LabelStmt:
		stmt, err = transpileLabelStmt(n, p)
		return

	case *ast.GotoStmt:
		stmt, err = transpileGotoStmt(n, p)
		return

	case *ast.GCCAsmStmt:
		// Go does not support inline assembly. See:
		// https://github.com/elliotchance/c2go/issues/228
		p.AddMessage(p.GenerateWarningMessage(
			errors.New("cannot transpile asm, will be ignored"), n))

		stmt = &goast.EmptyStmt{}
		return
	}

	// We do not care about the return type.
	expr, _, preStmts, postStmts, err = transpileToExpr(node, p, true)
	if err != nil {
		return
	}

	stmt = util.NewExprStmt(expr)

	return
}

func transpileToNode(node ast.Node, p *program.Program) error {
	switch n := node.(type) {
	case *ast.TranslationUnitDecl:
		for _, c := range n.Children() {
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
