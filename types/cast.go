package types

import (
	"fmt"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"
)

func Cast(program *program.Program, expr, fromType, toType string) string {
	fromType = ResolveType(program, fromType)
	toType = ResolveType(program, toType)

	if fromType == toType {
		return expr
	}

	types := []string{"int", "int64", "uint32", "__darwin_ct_rune_t",
		"byte", "float32", "float64"}

	for _, v := range types {
		if fromType == v && toType == "bool" {
			return fmt.Sprintf("%s != 0", expr)
		}
	}

	if fromType == "*int" && toType == "bool" {
		return fmt.Sprintf("%s != nil", expr)
	}

	if util.InStrings(fromType, types) && util.InStrings(toType, types) {
		return fmt.Sprintf("%s(%s)", toType, expr)
	}

	program.AddImport("github.com/elliotchance/c2go/noarch")
	return fmt.Sprintf("noarch.%sTo%s(%s)", util.Ucfirst(fromType), util.Ucfirst(toType), expr)
}
