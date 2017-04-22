package ast

import (
	"bytes"

	"github.com/elliotchance/c2go/program"
)

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

func (n *DeclStmt) render(program *program.Program) (string, string) {
	out := bytes.NewBuffer([]byte{})

	for _, child := range n.Children {
		src, _ := renderExpression(program, child)
		printLine(out, src, program.Indent)
	}

	return out.String(), ""
}

func (n *DeclStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
