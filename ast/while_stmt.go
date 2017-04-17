package ast

import (
	"bytes"
	"fmt"
)

type WhileStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseWhileStmt(line string) *WhileStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &WhileStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *WhileStmt) render(ast *Ast) (string, string) {
	out := bytes.NewBuffer([]byte{})
	// TODO: The first child of a WhileStmt appears to always be null.
	// Are there any cases where it is used?
	children := n.Children[1:]

	e, eType := renderExpression(ast, children[0])
	printLine(out, fmt.Sprintf("for %s {", cast(ast, e, eType, "bool")), ast.indent)

	body, _ := renderExpression(ast, children[1])
	printLine(out, body, ast.indent+1)

	printLine(out, "}", ast.indent)

	return out.String(), ""
}

func (n *WhileStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
