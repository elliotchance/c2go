package transpiler

import (
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"

	goast "go/ast"
)

func transpileDeclRefExpr(n *ast.DeclRefExpr, p *program.Program) (*goast.Ident, string, error) {
	if n.Name == "argc" {
		n.Name = "len(os.Args)"
		p.AddImport("os")
	} else if n.Name == "argv" {
		n.Name = "os.Args"
		p.AddImport("os")
	}

	return goast.NewIdent(n.Name), "", nil
}

func transpileDeclStmt(n *ast.DeclStmt, p *program.Program) (*goast.DeclStmt, error) {
	names := []*goast.Ident{}
	var theType string

	for _, c := range n.Children {
		if a, ok := c.(*ast.RecordDecl); ok {
			names = append(names, goast.NewIdent(a.Name))
		}
		if a, ok := c.(*ast.VarDecl); ok {
			names = append(names, goast.NewIdent(a.Name))
			theType = a.Type
		}
	}

	// panic(fmt.Sprintf("%#v", n.Children))

	return &goast.DeclStmt{
		Decl: &goast.GenDecl{
			Doc:    nil,
			TokPos: token.NoPos,
			Tok:    token.VAR,
			Lparen: token.NoPos,
			Specs: []goast.Spec{&goast.ValueSpec{
				Doc:     nil,
				Names:   names,
				Type:    goast.NewIdent(theType),
				Values:  nil,
				Comment: nil,
			}},
			Rparen: token.NoPos,
		},
	}, nil
}
