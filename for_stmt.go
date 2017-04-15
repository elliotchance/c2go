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
	// TODO: The second child of a ForStmt appears to always be null.
	// Are there any cases where it is used?
	if children[1] != nil {
		panic("non-nil child 1 in ForStmt")
	}
	b := renderExpression(children[2])[0]
	c := renderExpression(children[3])[0]

	if a == "" && b == "" && c == "" {
		printLine(out, fmt.Sprintf("for {"), indent)
	} else {
		printLine(out, fmt.Sprintf("for %s; %s; %s {", a, b, c), indent)
	}

	Render(out, children[4], functionName, indent+1, returnType)

	printLine(out, "}", indent)
}
