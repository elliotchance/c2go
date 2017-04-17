package ast

import (
	"bytes"
)

type ReturnStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseReturnStmt(line string) *ReturnStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ReturnStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *ReturnStmt) render(ast *Ast) (string, string) {
	out := bytes.NewBuffer([]byte{})
	r := "return"

	if len(n.Children) > 0 && ast.functionName != "main" {
		re, reType := renderExpression(ast, n.Children[0])
		r = "return " + cast(ast, re, reType, "int")
	}

	printLine(out, r, ast.indent)

	return out.String(), ""
}

func (n *ReturnStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
