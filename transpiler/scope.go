// This file contains functions for transpiling scopes. A scope is zero or more
// statements between a set of curly brackets.

package transpiler

import (
	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
)

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
		List: stmts,
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
