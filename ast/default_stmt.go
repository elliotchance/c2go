package ast

import (
	"github.com/elliotchance/c2go/program"
)

type DefaultStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseDefaultStmt(line string) *DefaultStmt {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &DefaultStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *DefaultStmt) render(program *program.Program) (string, string) {
	d := "default:"

	for _, s := range n.Children {
		line, _ := s.render(program)
		d += "\n" + line
	}

	return d, ""
}

func (n *DefaultStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
