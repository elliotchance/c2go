package transpiler

import (
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"

	goast "go/ast"
)

func transpileDeclStmt(n *ast.DeclStmt, p *program.Program) (*goast.DeclStmt, error) {
	names := []*goast.Ident{}

	for _, c := range n.Children {
		if a, ok := c.(*ast.RecordDecl); ok {
			names = append(names, goast.NewIdent(a.Name))
		}
		if a, ok := c.(*ast.VarDecl); ok {
			names = append(names, goast.NewIdent(a.Name))
		}
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
}
