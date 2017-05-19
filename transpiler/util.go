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
