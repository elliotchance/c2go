package ast

import (
	"strings"
)

// EnumDecl is node represents a enum declaration.
type EnumDecl struct {
	Addr       Address
	Pos        Position
	Position2  string
	Name       string
	ChildNodes []Node
}

func parseEnumDecl(line string) *EnumDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>(?P<position2> .+:\d+)?(?P<name>.*)`,
		line,
	)

	return &EnumDecl{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Position2:  groups["position2"],
		Name:       strings.TrimSpace(groups["name"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *EnumDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *EnumDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *EnumDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *EnumDecl) Position() Position {
	return n.Pos
}
