package ast

import "strings"

// IndirectFieldDecl is node represents a indirect field declaration.
type IndirectFieldDecl struct {
	Addr       Address
	Pos        Position
	Position2  string
	Implicit   bool
	Name       string
	Type       string
	ChildNodes []Node
}

func parseIndirectFieldDecl(line string) *IndirectFieldDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<position2> [^ ]+:[\d:]+)?
		(?P<implicit> implicit)?
		 (?P<name>\w+)
		 '(?P<type>.+?)'`,
		line,
	)

	return &IndirectFieldDecl{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Position2:  strings.TrimSpace(groups["position2"]),
		Implicit:   len(groups["implicit"]) > 0,
		Name:       groups["name"],
		Type:       groups["type"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *IndirectFieldDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *IndirectFieldDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *IndirectFieldDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *IndirectFieldDecl) Position() Position {
	return n.Pos
}
