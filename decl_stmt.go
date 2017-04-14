package main

import "bytes"

type DeclStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseDeclStmt(line string) *DeclStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &DeclStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *DeclStmt) RenderLine(out *bytes.Buffer, functionName string, indent int, returnType string) {
	for _, child := range n.Children {
		printLine(out, renderExpression(child)[0], indent)
	}
}
