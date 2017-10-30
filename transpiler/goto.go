// This file contains functions for transpiling goto/label statements.

package transpiler

import (
	goast "go/ast"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"
)

func transpileLabelStmt(n *ast.LabelStmt, p *program.Program) (*goast.LabeledStmt, []goast.Stmt, []goast.Stmt, error) {

	var post []goast.Stmt
	if len(n.Children()) > 0 {
		for _, node := range n.Children() {
			var stmt goast.Stmt
			stmt, preStmts, postStmts, err := transpileToStmt(node, p)
			if err != nil {
				return nil, nil, nil, err
			}
			post = append(post, preStmts...)
			// nil is happen, when transpile For
			if stmt != (*goast.ForStmt)(nil) {
				post = append(post, stmt)
			}
			post = append(post, postStmts...)
		}
	}

	return &goast.LabeledStmt{
		Label: util.NewIdent(n.Name),
		Stmt:  &goast.EmptyStmt{},
	}, []goast.Stmt{}, post, nil
}

func transpileGotoStmt(n *ast.GotoStmt, p *program.Program) (*goast.BranchStmt, error) {
	return &goast.BranchStmt{
		Label: util.NewIdent(n.Name),
		Tok:   token.GOTO,
	}, nil
}
