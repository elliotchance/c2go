package program

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

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
}

// NewProgram creates a new blank program.
func NewProgram() *Program {
	return &Program{
		imports:             []string{},
		typesAlreadyDefined: []string{},
	}
}

// Imports returns all of the Go imports for this program.
func (a *Program) Imports() []string {
	return a.imports
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

func (a *Program) ImportType(name string) string {
	if strings.Index(name, ".") != -1 {
		parts := strings.Split(name, ".")
		a.AddImport(strings.Join(parts[:len(parts)-1], "."))

		parts2 := strings.Split(name, "/")
		return parts2[len(parts2)-1]
	}

	return name
}

func (a *Program) TypeIsAlreadyDefined(typeName string) bool {
	return util.InStrings(typeName, a.typesAlreadyDefined)
}

func (a *Program) TypeIsNowDefined(typeName string) {
	a.typesAlreadyDefined = append(a.typesAlreadyDefined, typeName)
}
