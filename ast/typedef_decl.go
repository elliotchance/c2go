package ast

import (
	"strings"
)

// TypedefDecl is node represents a typedef declaration.
type TypedefDecl struct {
	Addr         Address
	Pos          Position
	Position2    string
	Name         string
	Type         string
	Type2        string
	IsImplicit   bool
	IsReferenced bool
	ChildNodes   []Node
}

func parseTypedefDecl(line string) *TypedefDecl {
	groups := groupsFromRegex(
		`(?:prev (?P<prev>0x[0-9a-f]+) )?
		<(?P<position><invalid sloc>|.*?)>
		(?P<position2> <invalid sloc>| col:\d+| line:\d+:\d+)?
		(?P<implicit> implicit)?
		(?P<referenced> referenced)?
		(?P<name> \w+)?
		(?P<type> '.*?')?
		(?P<type2>:'.*?')?`,
		line,
	)

	type2 := groups["type2"]
	if type2 != "" {
		type2 = type2[2 : len(type2)-1]
	}

	return &TypedefDecl{
		Addr:         ParseAddress(groups["address"]),
		Pos:          NewPositionFromString(groups["position"]),
		Position2:    strings.TrimSpace(groups["position2"]),
		Name:         strings.TrimSpace(groups["name"]),
		Type:         removeQuotes(groups["type"]),
		Type2:        type2,
		IsImplicit:   len(groups["implicit"]) > 0,
		IsReferenced: len(groups["referenced"]) > 0,
		ChildNodes:   []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *TypedefDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *TypedefDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *TypedefDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *TypedefDecl) Position() Position {
	return n.Pos
}
