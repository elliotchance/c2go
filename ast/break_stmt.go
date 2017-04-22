package ast

import (
	"github.com/elliotchance/c2go/program"
)

type BreakStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseBreakStmt(line string) *BreakStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &BreakStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *BreakStmt) render(program *program.Program) (string, string) {
	return "break", ""
}

func (n *BreakStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
