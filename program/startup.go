package program

import (
	goast "go/ast"
)

// AppendStartupStatement adds a new statement that must be executed when the
// program starts up before any other code. These are required to setup state
// for global variables like STDIN that might be referenced by the program.
func (p *Program) AppendStartupStatement(stmt goast.Stmt) {
	p.startupStatements = append(p.startupStatements, stmt)
}

// AppendStartupExpr is a convienience method for adding a new startup statement
// that is in the form of an expression.
func (p *Program) AppendStartupExpr(e goast.Expr) {
	p.AppendStartupStatement(&goast.ExprStmt{
		X: e,
	})
}

// StartupStatements returns the statements that will be executed before the
// program starts. These are required to setup state for global variables like
// STDIN that might be referenced by the program.
func (p *Program) StartupStatements() []goast.Stmt {
	return p.startupStatements
}
