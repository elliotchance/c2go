package ast

import (
	"strings"
)

// FunctionDecl is node represents a function declaration.
type FunctionDecl struct {
	Addr         Address
	Pos          Position
	Prev         string
	Parent       string
	Position2    string
	Name         string
	Type         string
	Type2        string
	IsExtern     bool
	IsImplicit   bool
	IsUsed       bool
	IsReferenced bool
	IsStatic     bool
	IsInline     bool
	ChildNodes   []Node
}

func parseFunctionDecl(line string) *FunctionDecl {
	groups := groupsFromRegex(
		`(?:parent (?P<parent>0x[0-9a-f]+) )?
		(?:prev (?P<prev>0x[0-9a-f]+) )?
		<(?P<position1>.*?)>
		(?P<position2> <scratch space>[^ ]+| [^ ]+)?
		(?P<implicit> implicit)?
		(?P<used> used)?
		(?P<referenced> referenced)?
		 (?P<name>[_\w]+)
		 '(?P<type>.*?)'
		(:'(?P<type2>.*?)')?
		(?P<extern> extern)?
		(?P<static> static)?
		(?P<inline> inline)?
		`,
		line,
	)

	return &FunctionDecl{
		Addr:         ParseAddress(groups["address"]),
		Pos:          NewPositionFromString(groups["position1"]),
		Parent:       groups["parent"],
		Prev:         groups["prev"],
		Position2:    strings.TrimSpace(groups["position2"]),
		Name:         groups["name"],
		Type:         groups["type"],
		Type2:        groups["type2"],
		IsExtern:     len(groups["extern"]) > 0,
		IsImplicit:   len(groups["implicit"]) > 0,
		IsUsed:       len(groups["used"]) > 0,
		IsReferenced: len(groups["referenced"]) > 0,
		IsStatic:     len(groups["static"]) > 0,
		IsInline:     len(groups["inline"]) > 0,
		ChildNodes:   []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FunctionDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *FunctionDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *FunctionDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *FunctionDecl) Position() Position {
	return n.Pos
}
