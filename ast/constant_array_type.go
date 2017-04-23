package ast

import (
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"
)

type ConstantArrayType struct {
	Address  string
	Type     string
	Size     int
	Children []Node
}

func parseConstantArrayType(line string) *ConstantArrayType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<size>\\d+)",
		line,
	)

	return &ConstantArrayType{
		Address:  groups["address"],
		Type:     groups["type"],
		Size:     util.Atoi(groups["size"]),
		Children: []Node{},
	}
}

func (n *ConstantArrayType) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *ConstantArrayType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
