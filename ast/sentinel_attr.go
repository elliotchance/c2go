package ast

import (
	"strings"

	"github.com/elliotchance/c2go/util"
)

// SentinelAttr is a type of attribute that is optionally attached to a variable
// or struct field definition.
type SentinelAttr struct {
	Addr       Address
	Pos        Position
	A          int
	B          int
	ChildNodes []Node
}

func parseSentinelAttr(line string) *SentinelAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>(?P<a> \d+)(?P<b> \d+)?`,
		line,
	)

	b := 0
	if groups["b"] != "" {
		b = util.Atoi(strings.TrimSpace(groups["b"]))
	}

	return &SentinelAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		A:          util.Atoi(strings.TrimSpace(groups["a"])),
		B:          b,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *SentinelAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *SentinelAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *SentinelAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *SentinelAttr) Position() Position {
	return n.Pos
}
