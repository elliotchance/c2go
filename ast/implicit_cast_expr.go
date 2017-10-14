package ast

import (
	"regexp"
	"strings"
)

type ImplicitCastExpr struct {
	Addr       Address
	Pos        Position
	Type       string
	Kind       string
	ChildNodes []Node
}

var regexImplicitCastExpr *regexp.Regexp

func init() {
	rx := "<(?P<position>.*)> '(?P<type>.*)' <(?P<kind>.*)>"
	fullRegexp := "(?P<address>[0-9a-fx]+) " +
		strings.Replace(strings.Replace(rx, "\n", "", -1), "\t", "", -1)
	regexImplicitCastExpr = regexp.MustCompile(fullRegexp)
}

func parseImplicitCastExpr(line string) *ImplicitCastExpr {
	groups := groupsFromRegex2(
		regexImplicitCastExpr,
		line,
	)

	return &ImplicitCastExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Kind:       groups["kind"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ImplicitCastExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ImplicitCastExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ImplicitCastExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ImplicitCastExpr) Position() Position {
	return n.Pos
}
