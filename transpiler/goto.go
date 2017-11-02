// This file contains functions for transpiling goto/label statements.

package transpiler

import (
	goast "go/ast"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"
)

func transpileLabelStmt(n *ast.LabelStmt, p *program.Program) (*goast.LabeledStmt, error) {
	var stmt goast.Stmt
	if len(n.Children()) > 0 {
		var err error
		stmt, _, _, err = transpileToStmt(n.Children()[0], p)
		if err != nil {
			return nil, err
		}
	}

	if stmt == nil {
		stmt = &goast.EmptyStmt{}
	}

	return &goast.LabeledStmt{
		Label: util.NewIdent(n.Name),
		Stmt:  stmt,
	}, nil
}

func transpileGotoStmt(n *ast.GotoStmt, p *program.Program) (*goast.BranchStmt, error) {
	return &goast.BranchStmt{
		Label: util.NewIdent(n.Name),
		Tok:   token.GOTO,
	}, nil
}
