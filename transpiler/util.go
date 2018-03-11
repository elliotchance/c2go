// This file contains utility and helper methods for the transpiler.

package transpiler

import (
	goast "go/ast"
	"reflect"
)

func isNil(stmt goast.Stmt) bool {
	if stmt == nil {
		return true
	}
	return reflect.ValueOf(stmt).IsNil()
}

func convertDeclToStmt(decls []goast.Decl) (stmts []goast.Stmt) {
	for i := range decls {
		if decls[i] != nil {
			stmts = append(stmts, &goast.DeclStmt{Decl: decls[i]})
		}
	}
	return
}

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
	out = make([]goast.Stmt, 0, len(stmts))
	for _, stmt := range stmts {
		if !isNil(stmt) {
			out = append(out, stmt)
		}
	}
	return
}

// combineStmts - combine elements to slice
func combineStmts(stmt goast.Stmt, preStmts, postStmts []goast.Stmt) (stmts []goast.Stmt) {
	stmts = make([]goast.Stmt, 0, 1+len(preStmts)+len(postStmts))

	preStmts = nilFilterStmts(preStmts)
	if preStmts != nil {
		stmts = append(stmts, preStmts...)
	}
	if !isNil(stmt) {
		stmts = append(stmts, stmt)
	}
	postStmts = nilFilterStmts(postStmts)
	if postStmts != nil {
		stmts = append(stmts, postStmts...)
	}
	return
}
