package program

import (
    "strconv"
    "strings"
)

// Imports returns all of the Go imports for this program.
func (p *Program) Imports() []string {
    imports := p.imports

    // Add "last-minute" import
    add_import := func(imp string) {
        imp = strconv.Quote(imp)
        for _, i := range imports {
            if i == imp {
                return
            }
        }

        imports = append(imports, imp)
    }

    // Imports packages for unions if at least an union is defined
    for _, t := range p.typesAlreadyDefined {
        _, ok := p.Unions[t]
        if ok {
            add_import("reflect")
            add_import("unsafe")
            break
        }
    }

    return imports
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
