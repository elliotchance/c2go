package ast

import (
	"strings"
)

type VarDecl struct {
	Addr         Address
	Pos          Position
	Position2    string
	Name         string
	Type         string
	Type2        string
	IsExtern     bool
	IsUsed       bool
	IsCInit      bool
	IsReferenced bool
	ChildNodes   []Node
}

func parseVarDecl(line string) *VarDecl {
	groups := groupsFromRegex(
		`(?:prev (?P<prev>0x[0-9a-f]+) )?
		<(?P<position>.*)>(?P<position2> .+:\d+)?
		(?P<used> used)?
		(?P<referenced> referenced)?
		(?P<name> \w+)?
		 '(?P<type>.+?)'
		(?P<type2>:'.*?')?
		(?P<extern> extern)?
		(?P<cinit> cinit)?`,
		line,
	)

	type2 := groups["type2"]
	if type2 != "" {
		type2 = type2[2 : len(type2)-1]
	}

	return &VarDecl{
		Addr:         ParseAddress(groups["address"]),
		Pos:          NewPositionFromString(groups["position"]),
		Position2:    strings.TrimSpace(groups["position2"]),
		Name:         strings.TrimSpace(groups["name"]),
		Type:         groups["type"],
		Type2:        type2,
		IsExtern:     len(groups["extern"]) > 0,
		IsUsed:       len(groups["used"]) > 0,
		IsCInit:      len(groups["cinit"]) > 0,
		IsReferenced: len(groups["referenced"]) > 0,
		ChildNodes:   []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *VarDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *VarDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *VarDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *VarDecl) Position() Position {
	return n.Pos
}
