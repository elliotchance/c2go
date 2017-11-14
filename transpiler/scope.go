// This file contains functions for transpiling scopes. A scope is zero or more
// statements between a set of curly brackets.

package transpiler

import (
	"fmt"
	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
)

func transpileCompoundStmt(n *ast.CompoundStmt, p *program.Program) (
	*goast.BlockStmt, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}
	stmts := []goast.Stmt{}

	for _, x := range n.Children() {
		result, err := transpileToStmts(x, p)
		if err != nil {
			return nil, nil, nil, err
		}

		if result != nil {
			stmts = append(stmts, result...)
		}
	}

	return &goast.BlockStmt{
		List: stmts,
	}, preStmts, postStmts, nil
}

func transpileToBlockStmt(node ast.Node, p *program.Program) (
	*goast.BlockStmt, []goast.Stmt, []goast.Stmt, error) {
	stmts, err := transpileToStmts(node, p)
	if err != nil {
		return nil, nil, nil, err
	}

	if len(stmts) == 1 {
		if block, ok := stmts[0].(*goast.BlockStmt); ok {
			return block, nil, nil, nil
		}
	}

	if stmts == nil {
		return nil, nil, nil, fmt.Errorf("Stmts inside Block cannot be nil")
	}

	return &goast.BlockStmt{
		List: stmts,
	}, nil, nil, nil
}
