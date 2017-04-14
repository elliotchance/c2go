package main

import "bytes"

type CompoundStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseCompoundStmt(line string) *CompoundStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &CompoundStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *CompoundStmt) RenderLine(out *bytes.Buffer, functionName string, indent int, returnType string) {
	for _, c := range n.Children {
		Render(out, c, functionName, indent, returnType)
	}
}
