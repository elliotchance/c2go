package program

import (
	"strconv"
	"strings"
)

// Imports returns all of the Go imports for this program.
func (p *Program) Imports() []string {
	return p.imports
}

// AddImport will append an absolute import if it is unique to the list of
// imports for this program.
func (p *Program) AddImport(importPath string) {
	quotedImportPath := strconv.Quote(importPath)

	if len(importPath) == 0 {
		return
	}

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

// ImportType imports a package for a fully qualified type and returns the local
// type name. For example:
//
//     t := p.ImportType("github.com/elliotchance/c2go/darwin.CtRuneT")
//
// Will import "github.com/elliotchance/c2go/darwin" and return (value of t)
// "darwin.CtRuneT".
func (p *Program) ImportType(name string) string {
	if strings.Contains(name, ".") {
		parts := strings.Split(name, ".")
		p.AddImport(strings.Join(parts[:len(parts)-1], "."))

		parts2 := strings.Split(name, "/")
		return parts2[len(parts2)-1]
	}

	return name
}
