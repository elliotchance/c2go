// This file contains functions for transpiling goto/label statements.

package transpiler

import (
	goast "go/ast"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"
)

func transpileLabelStmt(n *ast.LabelStmt, p *program.Program) (stmt *goast.LabeledStmt, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {

	if len(n.Children()) > 0 {
		var s goast.Stmt
		s, preStmts, postStmts, err = transpileToStmt(n.Children()[0], p)
		if err != nil {
			return nil, nil, nil, err
		}
		if s != (*goast.ForStmt)(nil) {
			var post []goast.Stmt
			post = append(post, s)
			postStmts = append(post, postStmts...)
		}
	}

	return &goast.LabeledStmt{
		Label: util.NewIdent(n.Name),
		Stmt:  &goast.EmptyStmt{},
	}, preStmts, postStmts, nil
}

func transpileGotoStmt(n *ast.GotoStmt, p *program.Program) (*goast.BranchStmt, error) {
	return &goast.BranchStmt{
		Label: util.NewIdent(n.Name),
		Tok:   token.GOTO,
	}, nil
}
