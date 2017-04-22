package program

import (
	"strings"
)

type Program struct {
	imports []string

	// for rendering go src
	FunctionName string
	Indent       int
	ReturnType   string
}

func NewProgram() *Program {
	return &Program{
		imports: []string{"fmt"},
	}
}

func (a *Program) Imports() []string {
	return a.imports
}

func (a *Program) AddImport(name string) {
	for _, i := range a.imports {
		if i == name {
			// already imported
			return
		}
	}

	a.imports = append(a.imports, name)
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
