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
	"github.com/elliotchance/c2go/types"
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
	decls, err := transpileToNode(root, p)
	if err != nil {
		p.AddMessage(p.GenerateErrorMessage(fmt.Errorf("Error of transpiling: err = %v", err), root))
		err = nil // Error is ignored
	}
	p.File.Decls = append(p.File.Decls, decls...)

	if p.OutputAsTest {
		p.AddImport("testing")
		p.AddImport("io/ioutil")
		p.AddImport("os")

		// TODO: There should be a cleaner way to add a function to the program.
		// This code was taken from the end of transpileFunctionDecl.
		p.File.Decls = append(p.File.Decls, &goast.FuncDecl{
			Name: util.NewIdent("TestApp"),
			Type: &goast.FuncType{
				Params: &goast.FieldList{
					List: []*goast.Field{
						{
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
		Name: goast.NewIdent("init"),
		Type: util.NewFuncType(&goast.FieldList{}, "", false),
		Body: &goast.BlockStmt{
			List: p.StartupStatements(),
		},
	})

	// only for "stdbool.h"
	if p.IncludeHeaderIsExists("stdbool.h") {
		p.File.Decls = append(p.File.Decls, &goast.GenDecl{
			Tok: token.TYPE,
			Specs: []goast.Spec{
				&goast.TypeSpec{
					Name: goast.NewIdent("_Bool"),
					Type: goast.NewIdent("int8"),
				},
			},
		})
	}

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
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileToExpr. err = %v", err)
		}
	}()
	if node == nil {
		panic(node)
	}
	defer func() {
		preStmts = nilFilterStmts(preStmts)
		postStmts = nilFilterStmts(postStmts)
	}()

	switch n := node.(type) {
	case *ast.StringLiteral:
		expr = transpileStringLiteral(n)
		exprType = "const char *"

	case *ast.FloatingLiteral:
		expr = transpileFloatingLiteral(n)
		exprType = "double"
		err = nil

	case *ast.PredefinedExpr:
		expr, exprType, err = transpilePredefinedExpr(n, p)

	case *ast.ConditionalOperator:
		expr, exprType, preStmts, postStmts, err = transpileConditionalOperator(n, p)

	case *ast.ArraySubscriptExpr:
		expr, exprType, preStmts, postStmts, err = transpileArraySubscriptExpr(n, p, exprIsStmt)

	case *ast.BinaryOperator:
		expr, exprType, preStmts, postStmts, err = transpileBinaryOperator(n, p, exprIsStmt)

	case *ast.UnaryOperator:
		expr, exprType, preStmts, postStmts, err = transpileUnaryOperator(n, p, exprIsStmt)

	case *ast.MemberExpr:
		expr, exprType, preStmts, postStmts, err = transpileMemberExpr(n, p)

	case *ast.ImplicitCastExpr:
		expr, exprType, preStmts, postStmts, err = transpileImplicitCastExpr(n, p, exprIsStmt)

	case *ast.DeclRefExpr:
		expr, exprType, err = transpileDeclRefExpr(n, p)

	case *ast.IntegerLiteral:
		expr, exprType, err = transpileIntegerLiteral(n), "int", nil

	case *ast.ParenExpr:
		expr, exprType, preStmts, postStmts, err = transpileParenExpr(n, p)

	case *ast.CStyleCastExpr:
		expr, exprType, preStmts, postStmts, err = transpileCStyleCastExpr(n, p, exprIsStmt)

	case *ast.CharacterLiteral:
		expr, exprType, err = transpileCharacterLiteral(n), "char", nil

	case *ast.CallExpr:
		expr, exprType, preStmts, postStmts, err = transpileCallExpr(n, p)

	case *ast.CompoundAssignOperator:
		return transpileCompoundAssignOperator(n, p, exprIsStmt)

	case *ast.UnaryExprOrTypeTraitExpr:
		return transpileUnaryExprOrTypeTraitExpr(n, p)

	case *ast.InitListExpr:
		expr, exprType, err = transpileInitListExpr(n, p)

	case *ast.CompoundLiteralExpr:
		expr, exprType, err = transpileCompoundLiteralExpr(n, p)

	case *ast.StmtExpr:
		return transpileStmtExpr(n, p)

	case *ast.ImplicitValueInitExpr:
		cType := n.Type1

		if strings.HasPrefix(cType, "struct ") {
			s := p.Structs[cType]
			if s == nil {
				return nil, "", nil, nil, fmt.Errorf("cannot found struct with name: `%s`", cType)
			}
			expr = &goast.CompositeLit{
				Type:   util.NewIdent(cType[len("struct "):]),
				Lbrace: 1,
			}
			return
		}

		s := p.Structs["struct "+cType]
		if s == nil {
			return nil, "", nil, nil, fmt.Errorf("cannot found struct with name: `%s`", cType)
		}
		expr = &goast.CompositeLit{
			Type:   util.NewIdent(cType),
			Lbrace: 1,
		}

	default:
		p.AddMessage(p.GenerateWarningMessage(errors.New("cannot transpile to expr"), node))
		expr = util.NewNil()
	}

	// Real return is through named arguments.
	return
}

func transpileToStmts(node ast.Node, p *program.Program) (stmts []goast.Stmt, err error) {
	if node == nil {
		return nil, nil
	}
	defer func() {
		stmts = nilFilterStmts(stmts)
	}()

	switch n := node.(type) {
	case *ast.DeclStmt:
		stmts, err = transpileDeclStmt(n, p)
		if err != nil {
			p.AddMessage(p.GenerateErrorMessage(fmt.Errorf("Error in DeclStmt: %v", err), n))
			err = nil // Error is ignored
		}
		return
	}

	var (
		stmt      goast.Stmt
		preStmts  []goast.Stmt
		postStmts []goast.Stmt
	)
	stmt, preStmts, postStmts, err = transpileToStmt(node, p)
	if err != nil {
		p.AddMessage(p.GenerateErrorMessage(fmt.Errorf("Error in DeclStmt: %v", err), node))
		err = nil // Error is ignored
	}
	return stripParentheses(combineStmts(stmt, preStmts, postStmts)), err
}

func stripParentheses(stmts []goast.Stmt) []goast.Stmt {
	for _, s := range stmts {
		if es, ok := s.(*goast.ExprStmt); ok {
			for {
				if pe, ok2 := es.X.(*goast.ParenExpr); ok2 {
					es.X = pe.X
				} else {
					break
				}
			}
		}
	}
	return stmts
}

func transpileToStmt(node ast.Node, p *program.Program) (
	stmt goast.Stmt, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	if node == nil {
		return
	}

	defer func() {
		if err != nil {
			p.AddMessage(p.GenerateErrorMessage(err, node))
			err = nil // Error is ignored
		}
	}()
	defer func() {
		preStmts = nilFilterStmts(preStmts)
		postStmts = nilFilterStmts(postStmts)
	}()

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
		stmt, preStmts, postStmts, err = transpileIfStmt(n, p)
		return

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
		stmt, preStmts, postStmts, err = transpileLabelStmt(n, p)
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
	case *ast.DeclStmt:
		var stmts []goast.Stmt
		stmts, err = transpileDeclStmt(n, p)
		if err != nil {
			return
		}
		stmt = stmts[len(stmts)-1]
		if len(stmts) > 1 {
			preStmts = stmts[0 : len(stmts)-2]
		}
		return
	}

	// We do not care about the return type.
	var theType string
	expr, theType, preStmts, postStmts, err = transpileToExpr(node, p, true)
	if err != nil {
		return
	}

	// nil is happen, when we remove function `free` of <stdlib.h>
	// see function CallExpr in transpiler
	if expr == (*goast.CallExpr)(nil) {
		return
	}

	// CStyleCastExpr.Kind == ToVoid
	var foundToVoid bool
	if theType == types.ToVoid {
		foundToVoid = true
	}
	if v, ok := node.(*ast.CStyleCastExpr); ok && v.Kind == ast.CStyleCastExprToVoid {
		foundToVoid = true
	}
	if len(node.Children()) > 0 {
		if v, ok := node.Children()[0].(*ast.CStyleCastExpr); ok && v.Kind == ast.CStyleCastExprToVoid {
			foundToVoid = true
		}
	}
	if foundToVoid {
		stmt = &goast.AssignStmt{
			Lhs: []goast.Expr{goast.NewIdent("_")},
			Tok: token.ASSIGN,
			Rhs: []goast.Expr{expr},
		}
		return
	}

	// For all other cases
	if expr == nil {
		err = fmt.Errorf("Expr is nil")
		return
	}
	stmt = util.NewExprStmt(expr)

	return
}

func transpileToNode(node ast.Node, p *program.Program) (decls []goast.Decl, err error) {
	defer func() {
		if err != nil {
			p.AddMessage(p.GenerateErrorMessage(err, node))
			err = nil // Error is ignored
		}
	}()

	switch n := node.(type) {
	case *ast.TranslationUnitDecl:
		decls, err = transpileTranslationUnitDecl(p, n)

	case *ast.FunctionDecl:
		decls, err = transpileFunctionDecl(n, p)
		if len(decls) > 0 {
			if _, ok := decls[0].(*goast.FuncDecl); ok {
				decls[0].(*goast.FuncDecl).Doc = p.GetMessageComments()
				decls[0].(*goast.FuncDecl).Doc.List =
					append(decls[0].(*goast.FuncDecl).Doc.List,
						p.GetComments(node.Position())...)
				decls[0].(*goast.FuncDecl).Doc.List =
					append([]*goast.Comment{&goast.Comment{
						Text: fmt.Sprintf("// %s - transpiled function from %s",
							decls[0].(*goast.FuncDecl).Name.Name,
							node.Position().GetSimpleLocation()),
					}}, decls[0].(*goast.FuncDecl).Doc.List...)
			}
		}

	case *ast.TypedefDecl:
		decls, err = transpileTypedefDecl(p, n)

	case *ast.RecordDecl:
		decls, err = transpileRecordDecl(p, n)

	case *ast.VarDecl:
		decls, _, err = transpileVarDecl(p, n)

	case *ast.EnumDecl:
		decls, err = transpileEnumDecl(p, n)

	case *ast.EmptyDecl:
		if len(n.Children()) == 0 {
			// ignore if length is zero, for avoid
			// mistake warning
		} else {
			p.AddMessage(p.GenerateWarningMessage(fmt.Errorf("EmptyDecl is not transpiled"), n))
		}
		err = nil
		return

	default:
		panic(fmt.Sprintf("cannot transpile to node: %#v", node))
	}

	return
}

func transpileStmts(nodes []ast.Node, p *program.Program) (stmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			p.AddMessage(p.GenerateErrorMessage(fmt.Errorf("Error in transpileToStmts: %v", err), nodes[0]))
			err = nil // Error is ignored
		}
	}()

	for _, s := range nodes {
		if s != nil {
			var (
				stmt      goast.Stmt
				preStmts  []goast.Stmt
				postStmts []goast.Stmt
			)
			stmt, preStmts, postStmts, err = transpileToStmt(s, p)
			if err != nil {
				return
			}
			stmts = append(stmts, combineStmts(stmt, preStmts, postStmts)...)
		}
	}

	return stmts, nil
}
