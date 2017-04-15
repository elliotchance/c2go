package main

import (
	"bytes"
	"fmt"
)

type IfStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseIfStmt(line string) *IfStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &IfStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *IfStmt) RenderLine(out *bytes.Buffer, functionName string, indent int, returnType string) {
	// TODO: The first two children of an IfStmt appear to always be null.
	// Are there any cases where they are used?
	children := n.Children[2:]

	e := renderExpression(children[0])
	printLine(out, fmt.Sprintf("if %s {", cast(e[0], e[1], "bool")), indent)

	Render(out, children[1], functionName, indent+1, returnType)

	if len(children) > 2 {
		printLine(out, "} else {", indent)
		Render(out, children[2], functionName, indent+1, returnType)
	}

	printLine(out, "}", indent)
}
