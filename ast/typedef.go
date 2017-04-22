package ast

import (
	"github.com/elliotchance/c2go/program"
)

type Typedef struct {
	Address  string
	Type     string
	Children []Node
}

func parseTypedef(line string) *Typedef {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return &Typedef{
		Address:  groups["address"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *Typedef) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *Typedef) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
