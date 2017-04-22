package ast

import (
	"github.com/elliotchance/c2go/program"
)

type ConstAttr struct {
	Address  string
	Position string
	Tags     string
	Children []Node
}

func parseConstAttr(line string) *ConstAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>(?P<tags>.*)",
		line,
	)

	return &ConstAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Tags:     groups["tags"],
		Children: []Node{},
	}
}

func (n *ConstAttr) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *ConstAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
