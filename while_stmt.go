package main

import (
	"bytes"
	"fmt"
)

type WhileStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseWhileStmt(line string) *WhileStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &WhileStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *WhileStmt) RenderLine(out *bytes.Buffer, functionName string, indent int, returnType string) {
	children := n.Children

	e := renderExpression(children[0])
	printLine(out, fmt.Sprintf("for %s {", cast(e[0], e[1], "bool")), indent)

	// FIXME: Does this do anything?
	Render(out, children[1], functionName, indent+1, returnType)

	printLine(out, "}", indent)
}
