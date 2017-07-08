package ast

import (
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

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ConstantArrayType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
