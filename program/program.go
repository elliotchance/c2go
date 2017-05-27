// Package program contains high-level orchestration and state of the input and
// output program during transpilation.
package program

import (
	"fmt"
	"go/token"

	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/util"
)

// Program contains all of the input, output and transpition state of a C
// program to a Go program.
type Program struct {
	// All of the Go import paths required for this program.
	imports []string

	// These are for the output Go AST.
	FileSet *token.FileSet
	File    *goast.File

	// One a type is defined it will be ignored if a future type of the same
	// name appears.
	typesAlreadyDefined []string

	// Contains the current function name during the transpilation.
	Function *ast.FunctionDecl

	// These are used to setup the runtime before the application begins. An
	// example would be to setup globals with stdin file pointers on certain
	// platforms.
	startupStatements []goast.Stmt

	// This is used to generate globally unique names for temporary variables
	// and other generated code. See GetNextIdentifier().
	nextUniqueIdentifier int

	// The definitions for defined structs.
	// TODO: This field should be protected through proper getters and setters.
	Structs map[string]*Struct

	// If verbose is on progress messages will be printed immediately as code
	// comments (so that they do not intefere with the program output).
	Verbose bool

	// A map of all the global variables (variables that exist outside of a
	// function) and their types.
	GlobalVariables map[string]string
}

// NewProgram creates a new blank program.
func NewProgram() *Program {
	return &Program{
		imports:             []string{},
		typesAlreadyDefined: []string{},
		startupStatements:   []goast.Stmt{},
		Structs:             make(map[string]*Struct),
		Verbose:             false,
		GlobalVariables:     map[string]string{},
	}
}

func (p *Program) TypeIsAlreadyDefined(typeName string) bool {
	return util.InStrings(typeName, p.typesAlreadyDefined)
}

func (p *Program) TypeIsNowDefined(typeName string) {
	p.typesAlreadyDefined = append(p.typesAlreadyDefined, typeName)
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
