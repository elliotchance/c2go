package ast

import (
	"regexp"
	"strings"
)

type BuiltinType struct {
	Addr       Address
	Type       string
	ChildNodes []Node
}

var regexBuiltinType *regexp.Regexp

func init() {
	rx := "'(?P<type>.*?)'"
	fullRegexp := "(?P<address>[0-9a-fx]+) " +
		strings.Replace(strings.Replace(rx, "\n", "", -1), "\t", "", -1)
	regexBuiltinType = regexp.MustCompile(fullRegexp)
}

func parseBuiltinType(line string) *BuiltinType {
	groups := groupsFromRegex2(
		regexBuiltinType,
		line,
	)

	return &BuiltinType{
		Addr:       ParseAddress(groups["address"]),
		Type:       groups["type"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *BuiltinType) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *BuiltinType) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *BuiltinType) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *BuiltinType) Position() Position {
	return Position{}
}
