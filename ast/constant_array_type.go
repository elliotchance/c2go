package ast

import (
	"github.com/elliotchance/c2go/util"
)

// ConstantArrayType is constant array type
type ConstantArrayType struct {
	Addr       Address
	Type       string
	Size       int
	ChildNodes []Node
}

func parseConstantArrayType(line string) *ConstantArrayType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<size>\\d+)",
		line,
	)

	return &ConstantArrayType{
		Addr:       ParseAddress(groups["address"]),
		Type:       groups["type"],
		Size:       util.Atoi(groups["size"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ConstantArrayType) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ConstantArrayType) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ConstantArrayType) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ConstantArrayType) Position() Position {
	return Position{}
}
