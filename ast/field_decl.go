package ast

import (
	"strings"
)

// FieldDecl is node represents a field declaration.
type FieldDecl struct {
	Addr       Address
	Pos        Position
	Position2  string
	Name       string
	Type       string
	Type2      string
	Implicit   bool
	Referenced bool
	ChildNodes []Node
}

func parseFieldDecl(line string) *FieldDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<position2> col:\d+| line:\d+:\d+)?
		(?P<implicit> implicit)?
		(?P<referenced> referenced)?
		(?P<name> \w+?)?
		 '(?P<type>.+?)'
		(:'(?P<type2>.*?)')?
		`,
		line,
	)

	return &FieldDecl{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Position2:  strings.TrimSpace(groups["position2"]),
		Name:       strings.TrimSpace(groups["name"]),
		Type:       groups["type"],
		Type2:      groups["type2"],
		Implicit:   len(groups["implicit"]) > 0,
		Referenced: len(groups["referenced"]) > 0,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FieldDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *FieldDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *FieldDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *FieldDecl) Position() Position {
	return n.Pos
}
