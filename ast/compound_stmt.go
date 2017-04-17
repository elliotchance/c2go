package ast

import "bytes"

type CompoundStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseCompoundStmt(line string) *CompoundStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &CompoundStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *CompoundStmt) render(ast *Ast) (string, string) {
	out := bytes.NewBuffer([]byte{})

	for _, c := range n.Children {
		src, _ := renderExpression(ast, c)
		printLine(out, src, ast.indent)
	}

	return out.String(), ""
}

func (n *CompoundStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
