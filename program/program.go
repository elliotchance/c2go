package program

import (
	"strings"

	"github.com/elliotchance/c2go/util"
)

type Program struct {
	imports []string

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

func (a *Program) TypeIsAlreadyDefined(typeName string) bool {
	return util.InStrings(typeName, a.TypesAlreadyDefined)
}

func (a *Program) TypeIsNowDefined(typeName string) {
	a.TypesAlreadyDefined = append(a.TypesAlreadyDefined, typeName)
}
