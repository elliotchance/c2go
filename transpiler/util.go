// This file contains utility and helper methods for the transpiler.

package transpiler

import (
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
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

// combineMultipleStmts - combine elements to slice
func combineMultipleStmts(stmts, preStmts, postStmts []goast.Stmt) []goast.Stmt {
	return combineStmts(nil, preStmts, append(stmts, postStmts...))
}

// GetUintptrForPointer - return uintptr for pointer
// Example : uint64(uintptr(unsafe.Pointer( ...pointer... )))
func GetUintptrForPointer(p *program.Program, expr goast.Expr, exprType string) (goast.Expr, string, error) {
	returnType := "long long"
	expr, err := types.CastExpr(p, expr, exprType, "void*")

	return util.NewCallExpr("int64",
		util.NewCallExpr("uintptr",
			expr)), returnType, err
}
