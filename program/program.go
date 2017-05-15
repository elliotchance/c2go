// Package program contains high-level orchestration and state of the input and
// output program during transpilation.
package program

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	goast "go/ast"

	"github.com/elliotchance/c2go/util"
)

// Program contains all of the input, output and transpition state of a C
// program to a Go program.
type Program struct {
	// All of the Go import paths required for this program.
	imports []string

	// These are for the output Go AST.
	FileSet *token.FileSet
	File    *ast.File

	// One a type is defined it will be ignored if a future type of the same
	// name appears.
	typesAlreadyDefined []string

	// Contains the current function name during the transpilation.
	FunctionName string

	// These are used to setup the runtime before the application begins. An
	// example would be to setup globals with stdin file pointers on certain
	// platforms.
	startupStatements []goast.Stmt

	// This is used to generate globally unique names for temporary variables
	// and other generated code. See GetNextIdentifier().
	nextUniqueIdentifier int

	// The definitions for defined structs.
	// TODO: This field should be protected through proper getters and setters.
	Structs map[string]Struct
}

// NewProgram creates a new blank program.
func NewProgram() *Program {
	return &Program{
		imports:             []string{},
		typesAlreadyDefined: []string{},
		startupStatements:   []goast.Stmt{},
		Structs:             make(map[string]Struct),
	}
}

// Imports returns all of the Go imports for this program.
func (p *Program) Imports() []string {
	return p.imports
}

// AddImport will append an absolute import if it is unique to the list of
// imports for this program.
func (p *Program) AddImport(importPath string) {
	quotedImportPath := strconv.Quote(importPath)

	for _, i := range p.imports {
		if i == quotedImportPath {
			// Already imported, ignore.
			return
		}
	}

	p.imports = append(p.imports, quotedImportPath)
}

// AddImports is a convienience method for adding multiple imports.
func (p *Program) AddImports(importPaths ...string) {
	for _, importPath := range importPaths {
		p.AddImport(importPath)
	}
}

func (p *Program) ImportType(name string) string {
	if strings.Index(name, ".") != -1 {
		parts := strings.Split(name, ".")
		p.AddImport(strings.Join(parts[:len(parts)-1], "."))

		parts2 := strings.Split(name, "/")
		return parts2[len(parts2)-1]
	}

	return name
}

func (p *Program) TypeIsAlreadyDefined(typeName string) bool {
	return util.InStrings(typeName, p.typesAlreadyDefined)
}

func (p *Program) TypeIsNowDefined(typeName string) {
	p.typesAlreadyDefined = append(p.typesAlreadyDefined, typeName)
}

func (p *Program) AppendStartupStatement(stmt goast.Stmt) {
	p.startupStatements = append(p.startupStatements, stmt)
}

func (p *Program) AppendStartupExpr(e goast.Expr) {
	p.AppendStartupStatement(&goast.ExprStmt{
		X: e,
	})
}

func (p *Program) StartupStatements() []goast.Stmt {
	return p.startupStatements
}

// GetNextIdentifier generates a new gloablly unique identifier name. This can
// be used for variables and functions in generated code.
//
// The value of prefix is only useful for readability in the code. If the prefix
// is an empty string then the prefix "__temp" will be used.
func (p *Program) GetNextIdentifier(prefix string) string {
	if prefix == "" {
		prefix = "temp"
	}

	identifierName := fmt.Sprintf("%s%d", prefix, p.nextUniqueIdentifier)
	p.nextUniqueIdentifier++

	return identifierName
}
