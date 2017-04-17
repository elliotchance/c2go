package ast

import "fmt"

func cast(ast *Ast, expr, fromType, toType string) string {
	fromType = resolveType(ast, fromType)
	toType = resolveType(ast, toType)

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

	if inStrings(fromType, types) && inStrings(toType, types) {
		return fmt.Sprintf("%s(%s)", toType, expr)
	}

	ast.addImport("github.com/elliotchance/c2go/noarch")
	return fmt.Sprintf("noarch.%sTo%s(%s)", ucfirst(fromType), ucfirst(toType), expr)
}
