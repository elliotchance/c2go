package ast

import (
	"strings"
)

type FieldDecl struct {
	Addr       Address
	Position   string
	Position2  string
	Name       string
	Type       string
	Referenced bool
	ChildNodes []Node
}

func parseFieldDecl(line string) *FieldDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<position2> col:\d+| line:\d+:\d+)?
		(?P<referenced> referenced)?
		(?P<name> \w+?)?
		 '(?P<type>.+?)'`,
		line,
	)

	return &FieldDecl{
		Addr:       ParseAddress(groups["address"]),
		Position:   groups["position"],
		Position2:  strings.TrimSpace(groups["position2"]),
		Name:       strings.TrimSpace(groups["name"]),
		Type:       groups["type"],
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
