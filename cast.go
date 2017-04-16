package main

import (
	"fmt"
	"regexp"
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
		"int", "int64", "uint16", "uint32", "byte",
		"float32", "float64",

		// Darwin specific:
		"__darwin_ct_rune_t", "darwin.Darwin_ct_rune_t",
	}
	for _, v := range types {
		if fromType == v && toType == "bool" {
			return fmt.Sprintf("%s != 0", expr)
		}
	}

	// In the forms of:
	// - `string` -> `[8]byte`
	// - `string` -> `char *[13]`
	match1 := regexp.MustCompile(`\[(\d+)\]byte`).FindStringSubmatch(toType)
	match2 := regexp.MustCompile(`char \*\[(\d+)\]`).FindStringSubmatch(toType)
	if fromType == "string" && (len(match1) > 0 || len(match2) > 0) {
		// Construct a byte array from "first":
		//
		//     var str [5]byte = [5]byte{'f','i','r','s','t'}

		s := ""
		for i := 1; i < len(expr)-1; i++ {
			if i > 1 {
				s += "','"
			}

			// Watch out for escape characters.
			if expr[i] == '\\' {
				s += fmt.Sprintf("\\%c", expr[i+1])
				i += 1
			} else {
				s += string(expr[i])
			}
		}

		size := "0"
		if len(match1) > 0 {
			size = match1[1]
		} else {
			size = match2[1]
		}

		return fmt.Sprintf("[%s]byte{'%s', 0}", size, s)
	}

	// In the forms of:
	// - `[7]byte` -> `string`
	// - `char *[12]` -> `string`
	match1 = regexp.MustCompile(`\[(\d+)\]byte`).FindStringSubmatch(fromType)
	match2 = regexp.MustCompile(`char \*\[(\d+)\]`).FindStringSubmatch(fromType)
	if (len(match1) > 0 || len(match2) > 0) && toType == "string" {
		size := 0
		if len(match1) > 0 {
			size = atoi(match1[1])
		} else {
			size = atoi(match2[1])
		}

		return fmt.Sprintf("string(%s[:%d])", expr, size-1)
	}

	// FIXME: This should be a more general rule for any pointer used a
	// bool.
	if fromType == "*int" && toType == "bool" {
		return fmt.Sprintf("%s != nil", expr)
	}
	if fromType == "int" && toType == "*int" {
		return "nil"
	}

	if inStrings(fromType, types) && inStrings(toType, types) {
		return fmt.Sprintf("%s(%s)", toType, expr)
	}

	addImport("github.com/elliotchance/c2go/noarch")
	return fmt.Sprintf("noarch.%sTo%s(%s)", ucfirst(fromType), ucfirst(toType), expr)
}
