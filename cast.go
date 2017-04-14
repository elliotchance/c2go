package main

import (
	"fmt"
	"strings"
)

func cast(expr, fromType, toType string) string {
	fromType = resolveType(fromType)
	toType = resolveType(toType)

	if fromType == toType {
		return expr
	}

	// Compatible integer types
	types := []string{
		// General types:
		"int", "int64", "uint32", "byte", "float32", "float64",

		// Darwin specific:
		"__darwin_ct_rune_t", "darwin.Darwin_ct_rune_t",
	}
	for _, v := range types {
		if fromType == v && toType == "bool" {
			return fmt.Sprintf("%s != 0", expr)
		}
	}

	if fromType == "string" && toType == "[8]byte" {
		return fmt.Sprintf("[8]byte{'%s'}",
			strings.Join(strings.Split(expr[1:len(expr)-1], ""), "','"))
	}

	if fromType == "*int" && toType == "bool" {
		return fmt.Sprintf("%s != nil", expr)
	}

	if inStrings(fromType, types) && inStrings(toType, types) {
		return fmt.Sprintf("%s(%s)", toType, expr)
	}

	addImport("github.com/elliotchance/c2go/noarch")
	return fmt.Sprintf("noarch.%sTo%s(%s)", ucfirst(fromType), ucfirst(toType), expr)
}
