package ast

import (
	"strings"

	"github.com/elliotchance/c2go/util"
)

// NonNullAttr is a type of attribute that is optionally attached to a variable
// or struct field definition.
type NonNullAttr struct {
	Addr       Address
	Pos        Position
	Inherited  bool
	A          int
	B          int
	C          int
	D          int
	ChildNodes []Node
}

func parseNonNullAttr(line string) *NonNullAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<inherited> Inherited)?
		(?P<a> \d+)?(?P<b> \d+)?(?P<c> \d+)?(?P<d> \d+)?`,
		line,
	)

	a := 0
	if groups["a"] != "" {
		a = util.Atoi(strings.TrimSpace(groups["a"]))
	}

	b := 0
	if groups["b"] != "" {
		b = util.Atoi(strings.TrimSpace(groups["b"]))
	}

	c := 0
	if groups["c"] != "" {
		c = util.Atoi(strings.TrimSpace(groups["c"]))
	}

	d := 0
	if groups["d"] != "" {
		d = util.Atoi(strings.TrimSpace(groups["d"]))
	}

	return &NonNullAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Inherited:  len(groups["inherited"]) > 0,
		A:          a,
		B:          b,
		C:          c,
		D:          d,
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
