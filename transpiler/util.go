// This file contains utility and helper methods for the transpiler.

package transpiler

import (
	"fmt"
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

func removeNil(nodes []goast.Decl) {
	var v V
	for _, node := range nodes {
		goast.Walk(v, node)
	}
}

type V struct{}

func (v V) Visit(node goast.Node) goast.Visitor {
	if node == (*goast.IfStmt)(nil) {
		fmt.Println("node  IfStmt - ", node)
	}
	if node == (*goast.ForStmt)(nil) {
		fmt.Println("node  ForStmt - ", node)
	}
	return v
}
