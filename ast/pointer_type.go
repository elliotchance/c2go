package ast

import (
	"github.com/elliotchance/c2go/program"
)

type PointerType struct {
	Address  string
	Type     string
	Children []Node
}

func parsePointerType(line string) *PointerType {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return &PointerType{
		Address:  groups["address"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *PointerType) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *PointerType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
