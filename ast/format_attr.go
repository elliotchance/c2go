package ast

import (
	"github.com/elliotchance/c2go/util"
)

// FormatAttr is a type of attribute that is optionally attached to a variable
// or struct field definition.
type FormatAttr struct {
	Addr         Address
	Pos          Position
	Implicit     bool
	Inherited    bool
	FunctionName string
	Unknown1     int
	Unknown2     int
	ChildNodes   []Node
}

func parseFormatAttr(line string) *FormatAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<implicit> Implicit)?
		(?P<inherited> Inherited)?
		 (?P<function>\w+)
		 (?P<unknown1>\d+)
		 (?P<unknown2>\d+)`,
		line,
	)

	return &FormatAttr{
		Addr:         ParseAddress(groups["address"]),
		Pos:          NewPositionFromString(groups["position"]),
		Implicit:     len(groups["implicit"]) > 0,
		Inherited:    len(groups["inherited"]) > 0,
		FunctionName: groups["function"],
		Unknown1:     util.Atoi(groups["unknown1"]),
		Unknown2:     util.Atoi(groups["unknown2"]),
		ChildNodes:   []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FormatAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *FormatAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *FormatAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *FormatAttr) Position() Position {
	return n.Pos
}
