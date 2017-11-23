package ast

import (
	"strings"

	"github.com/elliotchance/c2go/util"
)

type NonNullAttr struct {
	Addr       Address
	Pos        Position
	A          int
	B          int
	C          int
	ChildNodes []Node
}

func parseNonNullAttr(line string) *NonNullAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>(?P<a> \d+)(?P<b> \d+)?(?P<c> \d+)?`,
		line,
	)

	b := 0
	if groups["b"] != "" {
		b = util.Atoi(strings.TrimSpace(groups["b"]))
	}

	c := 0
	if groups["c"] != "" {
		c = util.Atoi(strings.TrimSpace(groups["c"]))
	}

	return &NonNullAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		A:          util.Atoi(strings.TrimSpace(groups["a"])),
		B:          b,
		C:          c,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *NonNullAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *NonNullAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *NonNullAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *NonNullAttr) Position() Position {
	return n.Pos
}
