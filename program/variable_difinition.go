package program

import "strings"

type VariableDefinition struct {
	// Headers "#include" with variable
	Headers []string

	// Name of variable in C source
	Cname string

	// Type of Go variable in C
	Ctype string

	// Name of variable in Go code
	GoName string

	// Type of Go variable GoName
	GoType string
}

func (v VariableDefinition) GetPackage() string {
	i := strings.Index(v.GoName, ".")
	if i < 0 {
		return ""
	}
	return "github.com/elliotchance/c2go/" + v.GoName[:i]
}

var builtInVariableDefinitions = []VariableDefinition{
	{
		Headers: []string{"unistd.h"},
		Cname:   "optarg",
		Ctype:   "char *",
		GoName:  "noarch.Optarg",
		GoType:  "[]byte",
	},
	{
		Headers: []string{"unistd.h"},
		Cname:   "opterr",
		Ctype:   "int",
		GoName:  "noarch.Opterr",
		GoType:  "int",
	},
	{
		Headers: []string{"unistd.h"},
		Cname:   "optind",
		Ctype:   "int",
		GoName:  "noarch.Optind",
		GoType:  "int",
	},
}

// IsBuiltInVariable - return true if that variable "name" is built-in
// variable of some C header
func (p *Program) IsBuiltInVariable(name string) bool {
	for i := range builtInVariableDefinitions {
		// TODO : add checking - if user C source have
		// included header
		if builtInVariableDefinitions[i].Cname == name {
			return true
		}
	}
	return false
}

func (p *Program) GetBuiltInVariableDefinition(name string) *VariableDefinition {
	for i := range builtInVariableDefinitions {
		// TODO : add checking - if user C source have
		// included header
		if builtInVariableDefinitions[i].Cname == name {
			return &builtInVariableDefinitions[i]
		}
	}
	return nil
}
