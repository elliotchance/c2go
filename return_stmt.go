package main

import "bytes"

type ReturnStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseReturnStmt(line string) *ReturnStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ReturnStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *ReturnStmt) RenderLine(out *bytes.Buffer, functionName string, indent int, returnType string) {
	r := "return"

	if len(n.Children) > 0 && functionName != "main" {
		re := renderExpression(n.Children[0])
		r = "return " + cast(re[0], re[1], "int")
	}

	printLine(out, r, indent)
}
