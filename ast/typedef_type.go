package ast

import (
	"github.com/elliotchance/c2go/program"
)

type TypedefType struct {
	Address  string
	Type     string
	Tags     string
	Children []Node
}

func parseTypedefType(line string) *TypedefType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<tags>.+)",
		line,
	)

	return &TypedefType{
		Address:  groups["address"],
		Type:     groups["type"],
		Tags:     groups["tags"],
		Children: []Node{},
	}
}

func (n *TypedefType) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *TypedefType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
