package main

import "strings"

var Imports []string

func init() {
	initImports()
}

func addImport(importName string) {
	for _, i := range Imports {
		if i == importName {
			return
		}
	}

	Imports = append(Imports, importName)
}

func importType(typeName string) string {
	if strings.Index(typeName, ".") != -1 {
		parts := strings.Split(typeName, ".")
		addImport(strings.Join(parts[:len(parts)-1], "."))

		parts2 := strings.Split(typeName, "/")
		return parts2[len(parts2)-1]
	}

	return typeName
}

func initImports() {
	Imports = []string{"fmt"}
}
