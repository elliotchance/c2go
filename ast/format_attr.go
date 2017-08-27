package ast

import (
	"github.com/elliotchance/c2go/util"
)

type FormatAttr struct {
	Addr         Address
	Position     string
	Implicit     bool
	Inherited    bool
	FunctionName string
	Unknown1     int
	Unknown2     int
	Children     []Node
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
		Position:     groups["position"],
		Implicit:     len(groups["implicit"]) > 0,
		Inherited:    len(groups["inherited"]) > 0,
		FunctionName: groups["function"],
		Unknown1:     util.Atoi(groups["unknown1"]),
		Unknown2:     util.Atoi(groups["unknown2"]),
		Children:     []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FormatAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *FormatAttr) Address() Address {
	return n.Addr
}
