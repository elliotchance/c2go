package main

import (
	"bytes"
)

type DoStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseDoStmt(line string) *DoStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &DoStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *DoStmt) RenderLine(out *bytes.Buffer, functionName string, indent int, returnType string) {
	// FIXME
	printLine(out, "/* do {} */", indent)
}
