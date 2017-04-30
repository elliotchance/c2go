package program

import (
	"go/ast"
	"go/token"
	"strings"

	"strconv"

	"github.com/elliotchance/c2go/util"
)

type Program struct {
	imports []string
	FileSet *token.FileSet
	File    *ast.File

	// for rendering go src
	TypesAlreadyDefined []string
	FunctionName        string
	Indent              int
	ReturnType          string
}

func NewProgram() *Program {
	return &Program{
		imports:             []string{},
		TypesAlreadyDefined: []string{},
	}
}

func (a *Program) Imports() []string {
	return a.imports
}

func (p *Program) AddImport(importPath string) {
	quotedImportPath := strconv.Quote(importPath)

	for _, i := range p.imports {
		if i == quotedImportPath {
			// already imported
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
	return util.InStrings(typeName, a.TypesAlreadyDefined)
}

func (a *Program) TypeIsNowDefined(typeName string) {
	a.TypesAlreadyDefined = append(a.TypesAlreadyDefined, typeName)
}
