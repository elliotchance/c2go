// This file contains utility and helper methods for the transpiler.

package transpiler

import (
	goast "go/ast"
)

func combinePreAndPostStmts(
	pre []goast.Stmt,
	post []goast.Stmt,
	newPre []goast.Stmt,
	newPost []goast.Stmt) ([]goast.Stmt, []goast.Stmt) {
	pre = append(pre, newPre...)
	post = append(post, newPost...)

	return pre, post
}

// combineStmts - combine elements to slice
func combineStmts(stmt goast.Stmt, preStmts, postStmts []goast.Stmt) (stmts []goast.Stmt) {
	if preStmts != nil {
		stmts = append(stmts, preStmts...)
	}
	if stmt != nil {
		stmts = append(stmts, stmt)
	}
	if postStmts != nil {
		stmts = append(stmts, postStmts...)
	}
	return
}
