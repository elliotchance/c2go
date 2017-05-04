package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
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

func transpileCompoundStmt(n *ast.CompoundStmt, p *program.Program) (*goast.BlockStmt, error) {
	stmts := []goast.Stmt{}

	for _, x := range n.Children {
		result, err := transpileToStmts(x, p)
		if err != nil {
			return nil, err
		}

		if result != nil {
			stmts = append(stmts, result...)
		}
	}

	return &goast.BlockStmt{
		Lbrace: token.NoPos,
		List:   stmts,
		Rbrace: token.NoPos,
	}, nil
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

func transpileFieldDecl(p *program.Program, n *ast.FieldDecl) (*goast.Field, string) {
	fieldType := types.ResolveType(p, n.Type)
	name := n.Name

	// FIXME: There are some cases where the name is empty. We need to
	// investigate this further. For now I will just exclude them.
	if name == "" {
		return nil, "unknown72"
	}

	// Go does not allow the name of a variable to be called "type". For the
	// moment I will rename this to avoid the error.
	if name == "type" {
		name = "type_"
	}

	// It may have a default value.
	// var defaultValue goast.Expr
	// if len(n.Children) > 0 {
	// 	defaultValue, _, err := transpileToExpr(n.Children[0], p)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// NULL is a macro that once rendered looks like "(0)" we have to be
	// sensitive to catch this as Go would complain that 0 (int) is not
	// compatible with the type we are setting it to.
	// if types.IsNullExpr(defaultValue) {
	// 	defaultValue = goast.NewIdent("nil")
	// }

	return &goast.Field{
		Names: []*goast.Ident{goast.NewIdent(name)},
		Type:  goast.NewIdent(fieldType),
	}, "unknown3"
}

func transpileRecordDecl(p *program.Program, n *ast.RecordDecl) error {
	name := n.Name
	if name == "" || p.TypeIsAlreadyDefined(name) {
		return nil
	}

	p.TypeIsNowDefined(name)

	if n.Kind == "union" {
		return nil
	}

	if name == "__locale_struct" ||
		name == "__sigaction" ||
		name == "sigaction" {
		return nil
	}

	var fields []*goast.Field
	for _, c := range n.Children {
		f, _ := transpileFieldDecl(p, c.(*ast.FieldDecl))
		fields = append(fields, f)
	}

	p.File.Decls = append(p.File.Decls, &goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name: goast.NewIdent(name),
				Type: &goast.StructType{
					Fields: &goast.FieldList{
						List: fields,
					},
				},
			},
		},
	})

	return nil
}

func transpileTypedefDecl(p *program.Program, n *ast.TypedefDecl) error {
	name := n.Name

	if p.TypeIsAlreadyDefined(name) {
		return nil
	}

	p.TypeIsNowDefined(name)

	resolvedType := types.ResolveType(p, n.Type)

	// There is a case where the name of the type is also the definition,
	// like:
	//
	//     type _RuneEntry _RuneEntry
	//
	// This of course is impossible and will cause the Go not to compile.
	// It itself is caused by lack of understanding (at this time) about
	// certain scenarios that types are defined as. The above example comes
	// from:
	//
	//     typedef struct {
	//        // ... some fields
	//     } _RuneEntry;
	//
	// Until which time that we actually need this to work I am going to
	// suppress these.
	if name == resolvedType {
		return nil
	}

	if name == "__mbstate_t" {
		resolvedType = p.ImportType("github.com/elliotchance/c2go/darwin.C__mbstate_t")
	}

	if name == "__darwin_ct_rune_t" {
		resolvedType = p.ImportType("github.com/elliotchance/c2go/darwin.Darwin_ct_rune_t")
	}

	// A bunch of random stuff to ignore... I really should deal with these.
	if name == "__builtin_va_list" ||
		name == "__qaddr_t" ||
		name == "definition" ||
		name == "_IO_lock_t" ||
		name == "va_list" ||
		name == "fpos_t" ||
		name == "__NSConstantString" ||
		name == "__darwin_va_list" ||
		name == "__fsid_t" ||
		name == "_G_fpos_t" ||
		name == "_G_fpos64_t" ||
		name == "__locale_t" ||
		name == "locale_t" ||
		name == "fsid_t" ||
		name == "sigset_t" {
		return nil
	}

	p.File.Decls = append(p.File.Decls, &goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name: goast.NewIdent(name),
				Type: goast.NewIdent(resolvedType),
			},
		},
	})

	return nil
}

func transpileVarDecl(p *program.Program, n *ast.VarDecl) string {
	theType := types.ResolveType(p, n.Type)
	name := n.Name

	// FIXME: These names don't seem to work when testing more than 1 file
	if name == "_LIB_VERSION" ||
		name == "_IO_2_1_stdin_" ||
		name == "_IO_2_1_stdout_" ||
		name == "_IO_2_1_stderr_" ||
		name == "stdin" ||
		name == "stdout" ||
		name == "stderr" ||
		name == "_DefaultRuneLocale" ||
		name == "_CurrentRuneLocale" {
		return "unknown10"
	}

	// Go does not allow the name of a variable to be called "type".
	// For the moment I will rename this to avoid the error.
	if name == "type" {
		name = "type_"
	}

	var defaultValues []goast.Expr
	if len(n.Children) > 0 {
		defaultValue, defaultValueType, err := transpileToExpr(n.Children[0], p)
		if err != nil {
			panic(err)
		}

		defaultValues = []goast.Expr{
			types.CastExpr(p, defaultValue, defaultValueType, n.Type),
			// &goast.BasicLit{
			// 	// TODO: It this safe to always be a STRING?
			// 	Kind:  token.STRING,
			// 	Value: types.CastExpr(p, defaultValue, defaultValueType, n.Type),
			// },
		}
	}

	// if suffix == " = (0)" {
	// 	suffix = " = nil"
	// }

	p.File.Decls = append(p.File.Decls, &goast.GenDecl{
		Tok: token.VAR,
		Specs: []goast.Spec{
			&goast.ValueSpec{
				Names: []*goast.Ident{
					goast.NewIdent(name),
				},
				Type:   goast.NewIdent(theType),
				Values: defaultValues,
			},
		},
	})

	return n.Type
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
