package main

import "bytes"

type TranslationUnitDecl struct {
	Address  string
	Children []interface{}
}

func parseTranslationUnitDecl(line string) *TranslationUnitDecl {
	groups := groupsFromRegex("", line)

	return &TranslationUnitDecl{
		Address:  groups["address"],
		Children: []interface{}{},
	}
}

func (n *TranslationUnitDecl) RenderLine(out *bytes.Buffer, functionName string, indent int, returnType string) {
	for _, c := range n.Children {
		Render(out, c, functionName, indent, returnType)
	}
}
