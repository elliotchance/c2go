package transpiler

import (
	"fmt"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"

	goast "go/ast"
)

func transpileDeclRefExpr(n *ast.DeclRefExpr, p *program.Program) (*goast.Ident, string, error) {
	// TODO: System arguments are fixed variable names.
	// https://github.com/elliotchance/c2go/issues/86
	if n.Name == "argc" {
		n.Name = "len(os.Args)"
		p.AddImport("os")
	}
	if n.Name == "argv" {
		n.Name = "os.Args"
		p.AddImport("os")
	}

	return goast.NewIdent(n.Name), n.Type, nil
}

func newDeclStmt(a *ast.VarDecl, p *program.Program) (*goast.DeclStmt, error) {
	var values []goast.Expr = nil
	if len(a.Children) > 0 {
		defaultValue, defaultValueType, err := transpileToExpr(a.Children[0], p)
		if err != nil {
			return nil, err
		}

		if !types.IsNullExpr(defaultValue) {
			values = []goast.Expr{
				types.CastExpr(p, defaultValue, defaultValueType, a.Type),
			}
		}
	}

	return &goast.DeclStmt{
		Decl: &goast.GenDecl{
			Tok: token.VAR,
			Specs: []goast.Spec{
				&goast.ValueSpec{
					Names:  []*goast.Ident{goast.NewIdent(a.Name)},
					Type:   goast.NewIdent(types.ResolveType(p, a.Type)),
					Values: values,
				},
			},
		},
	}, nil
}

func transpileDeclStmt(n *ast.DeclStmt, p *program.Program) ([]goast.Stmt, error) {
	// There may be more than one variable defined on the same line. With C it
	// is possible for them to have similar by diffrent types, whereas in Go
	// this is not possible. The easiest way around this is to split the
	// variables up into multiple declarations. That is why this function
	// returns one or more DeclStmts.
	decls := []goast.Stmt{}

	for _, c := range n.Children {
		switch a := c.(type) {
		case *ast.RecordDecl:
			// I'm not sure why this is ignored. Maybe we haven't found a
			// situation where this is needed yet?

		case *ast.VarDecl:
			e, err := newDeclStmt(a, p)
			if err != nil {
				return nil, err
			}

			decls = append(decls, e)

		default:
			panic(a)
		}
	}

	return decls, nil
}

func transpileArraySubscriptExpr(n *ast.ArraySubscriptExpr, p *program.Program) (*goast.IndexExpr, string, error) {
	children := n.Children
	expression, expressionType, err := transpileToExpr(children[0], p)
	if err != nil {
		return nil, "", err
	}

	index, _, err := transpileToExpr(children[1], p)
	if err != nil {
		return nil, "", err
	}

	newType, err := types.GetDereferenceType(expressionType)
	if err != nil {
		panic(fmt.Sprintf("Cannot dereference type '%s' for the expression '%s'",
			expressionType, expression))
	}

	return &goast.IndexExpr{
		X:     expression,
		Index: index,
	}, newType, nil
}

func transpileMemberExpr(n *ast.MemberExpr, p *program.Program) (*goast.SelectorExpr, string, error) {
	lhs, lhsType, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", err
	}

	lhsResolvedType := types.ResolveType(p, lhsType)
	rhs := n.Name

	// FIXME: This should not be empty. We need some fallback type to catch
	// errors like "unknown8".
	rhsType := ""

	// FIXME: This is just a hack
	if util.InStrings(lhsResolvedType, []string{"darwin.Float2", "darwin.Double2"}) {
		rhs = util.GetExportedName(rhs)
		rhsType = "int"
	}

	return &goast.SelectorExpr{
		X:   lhs,
		Sel: goast.NewIdent(rhs),
	}, rhsType, nil
}
