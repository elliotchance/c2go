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
	pre = append(pre, nilFilterStmts(newPre)...)
	post = append(post, nilFilterStmts(newPost)...)

	return pre, post
}

// nilFilterStmts - remove nil stmt from slice
func nilFilterStmts(stmts []goast.Stmt) (out []goast.Stmt) {
	for _, stmt := range stmts {
		if stmt != nil && stmt != (*goast.IfStmt)(nil) && stmt != (goast.Stmt)(nil) {
			out = append(out, stmt)
		}
	}
	return
}

// combineStmts - combine elements to slice
func combineStmts(stmt goast.Stmt, preStmts, postStmts []goast.Stmt) (stmts []goast.Stmt) {
	preStmts = nilFilterStmts(preStmts)
	if preStmts != nil {
		stmts = append(stmts, preStmts...)
	}
	if stmt != nil {
		stmts = append(stmts, stmt)
	}
	postStmts = nilFilterStmts(postStmts)
	if postStmts != nil {
		stmts = append(stmts, postStmts...)
	}
	return
}
