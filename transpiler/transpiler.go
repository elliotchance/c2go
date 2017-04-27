package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"strings"

	"strconv"

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
	return err
}

func transpileToExpr(node ast.Node, p *program.Program) (goast.Expr, string, error) {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return transpileStringLiteral(n), "", nil

	case *ast.PredefinedExpr:
		if n.Name == "__PRETTY_FUNCTION__" {
			// FIXME
			return &goast.BasicLit{
				Kind:  token.STRING,
				Value: "\"void print_number(int *)\"",
			}, "const char*", nil
		}

		if n.Name == "__func__" {
			// FIXME
			src := fmt.Sprintf("\"%s\"", "print_number")
			return &goast.BasicLit{
				Kind:  token.STRING,
				Value: src,
			}, "const char*", nil
		}

		panic(fmt.Sprintf("renderExpression: unknown PredefinedExpr: %s", n.Name))

	case *ast.ConditionalOperator:
		// TODO: check errors for these
		a, _, _ := transpileToExpr(n.Children[0], p)
		b, _, _ := transpileToExpr(n.Children[1], p)
		c, _, _ := transpileToExpr(n.Children[2], p)

		p.AddImport("github.com/elliotchance/c2go/noarch")

		return &goast.CallExpr{
			Fun:      goast.NewIdent("noarch.Ternary"),
			Lparen:   token.NoPos,
			Args:     []goast.Expr{a, b, c},
			Ellipsis: token.NoPos,
			Rparen:   token.NoPos,
		}, "", nil

		// src := fmt.Sprintf("noarch.Ternary(%s, func () interface{} { return %s }, func () interface{} { return %s })", a, b, c)
		// return src, n.Type

	case *ast.ArraySubscriptExpr:
		children := n.Children
		expression, expressionType, err := transpileToExpr(children[0], p)
		if err != nil {
			return nil, "", err
		}

		index, _, err := transpileToExpr(children[1], p)
		if err != nil {
			return nil, "", err
		}

		return &goast.IndexExpr{
			X:     expression,
			Index: index,
		}, expressionType, nil

	case *ast.BinaryOperator:
		left, _, err := transpileToExpr(n.Children[0], p)
		if err != nil {
			return nil, "", err
		}

		right, _, err := transpileToExpr(n.Children[1], p)
		if err != nil {
			return nil, "", err
		}

		return &goast.BinaryExpr{
			X:     left,
			OpPos: token.NoPos,
			Op:    getTokenForOperator(n.Operator),
			Y:     right,
		}, "", nil

	case *ast.UnaryOperator:
		left, _, err := transpileToExpr(n.Children[0], p)
		if err != nil {
			return nil, "", err
		}

		return &goast.UnaryExpr{
			X:     left,
			OpPos: token.NoPos,
			Op:    getTokenForOperator(n.Operator),
		}, "", nil

	case *ast.MemberExpr:
		lhs, _, err := transpileToExpr(n.Children[0], p)
		if err != nil {
			return nil, "", err
		}

		// lhsResolvedType := types.ResolveType(program, lhsType)
		rhs := n.Name
		// rhsType := ""

		// FIXME: This is just a hack
		// if util.InStrings(lhsResolvedType, []string{"darwin.Float2", "darwin.Double2"}) {
		// 	rhs = util.GetExportedName(rhs)
		// 	rhsType = "int"
		// }

		return &goast.SelectorExpr{
			X:   lhs,
			Sel: goast.NewIdent(rhs),
		}, "", nil

	case *ast.ImplicitCastExpr:
		return transpileToExpr(n.Children[0], p)

	case *ast.DeclRefExpr:
		return goast.NewIdent(n.Name), "", nil

	case *ast.IntegerLiteral:
		return &goast.BasicLit{
			ValuePos: token.NoPos,
			Kind:     token.INT,
			Value:    strconv.Itoa(n.Value),
		}, "", nil

	case *ast.ParenExpr:
		e, _, err := transpileToExpr(n.Children[0], p)
		if err != nil {
			return nil, "", err
		}

		return &goast.ParenExpr{
			Lparen: token.NoPos,
			X:      e,
			Rparen: token.NoPos,
		}, "", nil

	case *ast.CStyleCastExpr:
		return transpileToExpr(n.Children[0], p)

	case *ast.CharacterLiteral:
		return &goast.BasicLit{
			ValuePos: token.NoPos,
			Kind:     token.CHAR,
			Value:    fmt.Sprintf("%c", n.Value),
		}, "", nil

	case *ast.CallExpr:
		functionName := n.Children[0].(*ast.ImplicitCastExpr).Children[0].(*ast.DeclRefExpr).Name
		functionDef := program.GetFunctionDefinition(functionName)

		if functionDef == nil {
			panic(fmt.Sprintf("unknown function: %s", functionName))
		}

		if functionDef.Substitution != "" {
			parts := strings.Split(functionDef.Substitution, ".")
			importName := strings.Join(parts[:len(parts)-1], ".")
			p.AddImport(importName)

			parts2 := strings.Split(functionDef.Substitution, "/")
			functionName = parts2[len(parts2)-1]
		}

		args := []goast.Expr{}
		i := 0
		for _, arg := range n.Children[1:] {
			e, eType, err := transpileToExpr(arg, p)
			if err != nil {
				return nil, "", err
			}

			if i > len(functionDef.ArgumentTypes)-1 {
				// This means the argument is one of the varargs
				// so we don't know what type it needs to be
				// cast to.
				args = append(args, e)
			} else {
				args = append(args, types.CastExpr(p, e, eType, functionDef.ArgumentTypes[i]))
			}

			i++
		}

		return &goast.CallExpr{
			Fun:      goast.NewIdent(functionName),
			Lparen:   token.NoPos,
			Args:     args,
			Ellipsis: token.NoPos,
			Rparen:   token.NoPos,
		}, "", nil

		// src := fmt.Sprintf("%s(%s)", functionName, strings.Join(parts, ", "))
		// return src, functionDef.ReturnType
	}

	panic(fmt.Sprintf("cannot transpile to expr: %#v", node))
}

func transpileToStmt(node ast.Node, p *program.Program) (goast.Stmt, error) {
	switch n := node.(type) {
	case *ast.IfStmt:
		children := n.Children

		// There is always 4 or 5 children in an IfStmt. For example:
		//
		//     if (i == 0) {
		//         return 0;
		//     } else {
		//         return 1;
		//     }
		//
		// 1. Not sure what this is for. This gets removed.
		// 2. Not sure what this is for.
		// 3. conditional = BinaryOperator: i == 0
		// 4. body = CompoundStmt: { return 0; }
		// 5. elseBody = CompoundStmt: { return 1; }
		//
		// elseBody will be nil if there is no else clause.

		// On linux I have seen only 4 children for an IfStmt with the same
		// definitions above, but missing the first argument. Since we don't
		// know what the first argument is for anyway we will just remove it on
		// Mac if necessary.
		if len(children) == 5 && children[0] != nil {
			panic("non-nil child 0 in IfStmt")
		}
		if len(children) == 5 {
			children = children[1:]
		}

		// From here on there must be 4 children.
		if len(children) != 4 {
			panic(fmt.Sprintf("Expected 4 children in IfStmt, got %#v", children))
		}

		// Maybe we will discover what the nil value is?
		if children[0] != nil {
			panic("non-nil child 0 in IfStmt")
		}

		conditional, conditionalType, err := transpileToExpr(children[1], p)
		if err != nil {
			return nil, err
		}

		// The condition in Go must always be a bool.
		boolCondition := types.CastExpr(p, conditional, conditionalType, "bool")

		body, err := transpileToBlockStmt(children[2], p)
		if err != nil {
			return nil, err
		}

		r := &goast.IfStmt{
			If:   token.NoPos,
			Init: nil,
			Cond: boolCondition,
			Body: body,
		}

		if children[3] != nil {
			elseBody, err := transpileToBlockStmt(children[3], p)
			if err != nil {
				return nil, err
			}

			r.Else = elseBody
		}

		return r, nil

	case *ast.ForStmt:
		children := n.Children

		// There are always 5 children in a ForStmt, for example:
		//
		//     for ( c = 0 ; c < n ; c++ ) {
		//         doSomething();
		//     }
		//
		// 1. initExpression = BinaryStmt: c = 0
		// 2. Not sure what this is for, but it's always nil. There is a panic
		//    below in case we discover what it is used for (pun intended).
		// 3. conditionalExpression = BinaryStmt: c < n
		// 4. stepExpression = BinaryStmt: c++
		// 5. body = CompoundStmt: { CallExpr }

		if len(children) != 5 {
			panic(fmt.Sprintf("Expected 5 children in ForStmt, got %#v", children))
		}

		// TODO: The second child of a ForStmt appears to always be null.
		// Are there any cases where it is used?
		if children[1] != nil {
			panic("non-nil child 1 in ForStmt")
		}

		init, _ := transpileToStmt(children[0], p)
		conditional, _, _ := transpileToExpr(children[2], p)
		post, _ := transpileToStmt(children[3], p)
		body, _ := transpileToBlockStmt(children[4], p)

		return &goast.ForStmt{
			Init: init,
			Cond: conditional,
			Post: post,
			Body: body,
		}, nil

	case *ast.DeclStmt:
		names := []*goast.Ident{}

		for _, c := range n.Children {
			names = append(names, goast.NewIdent(c.(*ast.VarDecl).Name))
		}

		return &goast.DeclStmt{
			Decl: &goast.GenDecl{
				Doc:    nil,
				TokPos: token.NoPos,
				Tok:    token.VAR,
				Lparen: token.NoPos,
				Specs: []goast.Spec{&goast.ValueSpec{
					Doc:     nil,
					Names:   names,
					Type:    nil,
					Values:  nil,
					Comment: nil,
				}},
				Rparen: token.NoPos,
			},
		}, nil

	case *ast.ReturnStmt:
		return &goast.ReturnStmt{
			Return:  token.NoPos,
			Results: nil,
		}, nil

	case *ast.BreakStmt:
		return &goast.BranchStmt{
			Tok: token.BREAK,
		}, nil

	case *ast.CompoundStmt:
		stmts := []goast.Stmt{}

		for _, x := range n.Children {
			result, err := transpileToStmt(x, p)
			if err != nil {
				return nil, err
			}

			stmts = append(stmts, result)
		}

		return &goast.BlockStmt{
			Lbrace: token.NoPos,
			List:   stmts,
			Rbrace: token.NoPos,
		}, nil
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
		var body *goast.BlockStmt

		// Always register the new function. Only from this point onwards will
		// we be allowed to refer to the function.
		if program.GetFunctionDefinition(n.Name) == nil {
			program.AddFunctionDefinition(program.FunctionDefinition{
				Name:       n.Name,
				ReturnType: "int",
				// FIXME
				ArgumentTypes: []string{},
				Substitution:  "",
			})
		}

		// If the function has a direct substitute in Go we do not want to
		// output the C definition of it.
		if f := program.GetFunctionDefinition(n.Name); f != nil &&
			f.Substitution != "" {
			return nil
		}

		hasBody := false
		for _, c := range n.Children {
			if b, ok := c.(*ast.CompoundStmt); ok {
				var err error
				body, err = transpileToBlockStmt(b, p)
				if err != nil {
					return err
				}

				hasBody = true
				break
			}
		}

		if n.Name == "__istype" ||
			n.Name == "__isctype" ||
			n.Name == "__wcwidth" ||
			n.Name == "__sputc" ||
			n.Name == "__inline_signbitf" ||
			n.Name == "__inline_signbitd" ||
			n.Name == "__inline_signbitl" {
			return nil
		}

		if hasBody {
			p.File.Decls = append(p.File.Decls, &goast.FuncDecl{
				Doc:  nil,
				Recv: nil,
				Name: goast.NewIdent(n.Name),
				Type: &goast.FuncType{
					Params:  getFieldList(n),
					Results: nil,
				},
				Body: body,
			})
		}

	case *ast.TypedefDecl, *ast.RecordDecl, *ast.VarDecl:
		// do nothing

	default:
		panic(fmt.Sprintf("cannot transpile to node: %#v", node))
	}

	return nil
}

func getFieldList(f *ast.FunctionDecl) *goast.FieldList {
	r := []*goast.Field{}
	for _, n := range f.Children {
		if v, ok := n.(*ast.ParmVarDecl); ok {
			r = append(r, &goast.Field{
				Doc:     nil,
				Names:   []*goast.Ident{goast.NewIdent(v.Name)},
				Type:    goast.NewIdent(v.Type),
				Tag:     nil,
				Comment: nil,
			})
		}
	}

	return &goast.FieldList{Opening: token.NoPos, List: r, Closing: token.NoPos}
}
