package ast

import (
	"strings"
)

// ParmVarDecl is node represents a parameter of variable declaration.
type ParmVarDecl struct {
	Addr         Address
	Pos          Position
	Position2    string
	Name         string
	Type         string
	Type2        string
	IsUsed       bool
	IsReferenced bool
	IsRegister   bool
	ChildNodes   []Node
}

func parseParmVarDecl(line string) *ParmVarDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<position2> [^ ]+:[\d:]+)?
		(?P<used> used)?
		(?P<referenced> referenced)?
		(?P<name> \w+)?
		 '(?P<type>.*?)'
		(?P<type2>:'.*?')?
		(?P<register> register)?
		`,
		line,
	)

	type2 := groups["type2"]
	if type2 != "" {
		type2 = type2[2 : len(type2)-1]
	}

	if strings.Index(groups["position"], "<invalid sloc>") > -1 {
		groups["position"] = "<invalid sloc>"
		groups["position2"] = "<invalid sloc>"
	}

	return &ParmVarDecl{
		Addr:         ParseAddress(groups["address"]),
		Pos:          NewPositionFromString(groups["position"]),
		Position2:    strings.TrimSpace(groups["position2"]),
		Name:         strings.TrimSpace(groups["name"]),
		Type:         groups["type"],
		Type2:        type2,
		IsUsed:       len(groups["used"]) > 0,
		IsReferenced: len(groups["referenced"]) > 0,
		IsRegister:   len(groups["register"]) > 0,
		ChildNodes:   []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ParmVarDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ParmVarDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ParmVarDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ParmVarDecl) Position() Position {
	return n.Pos
}
