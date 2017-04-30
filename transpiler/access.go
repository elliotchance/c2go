package transpiler

import (
	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"

	goast "go/ast"
)

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

	return &goast.IndexExpr{
		X:     expression,
		Index: index,
	}, expressionType, nil
}

func transpileMemberExpr(n *ast.MemberExpr, p *program.Program) (*goast.SelectorExpr, string, error) {
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
}
