package main

import (
	"bytes"
	"fmt"
)

type ForStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseForStmt(line string) *ForStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ForStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *ForStmt) RenderLine(out *bytes.Buffer, functionName string, indent int, returnType string) {
	children := n.Children

	a := renderExpression(children[0])[0]
	b := renderExpression(children[1])[0]
	c := renderExpression(children[2])[0]

	printLine(out, fmt.Sprintf("for %s; %s; %s {", a, b, c), indent)

	Render(out, children[3], functionName, indent+1, returnType)

	printLine(out, "}", indent)
}
