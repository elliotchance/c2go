package ast

import "bytes"

type DeclStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseDeclStmt(line string) *DeclStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &DeclStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *DeclStmt) render(ast *Ast) (string, string) {
	out := bytes.NewBuffer([]byte{})

	for _, child := range n.Children {
		src, _ := renderExpression(ast, child)
		printLine(out, src, ast.indent)
	}

	return out.String(), ""
}

func (n *DeclStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
