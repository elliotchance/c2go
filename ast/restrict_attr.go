package ast

import (
	"github.com/elliotchance/c2go/program"
)

type RestrictAttr struct {
	Address  string
	Position string
	Name     string
	Children []Node
}

func parseRestrictAttr(line string) *RestrictAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> (?P<name>.+)",
		line,
	)

	return &RestrictAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Name:     groups["name"],
		Children: []Node{},
	}
}

func (n *RestrictAttr) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *RestrictAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
