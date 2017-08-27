package ast

import (
	"github.com/elliotchance/c2go/util"
)

type ConstantArrayType struct {
	Addr     Address
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
		Addr:     ParseAddress(groups["address"]),
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

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ConstantArrayType) Address() Address {
	return n.Addr
}
