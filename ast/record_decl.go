package ast

import (
	"strings"
)

// RecordDecl is node represents a record declaration.
type RecordDecl struct {
	Addr       Address
	Pos        Position
	Prev       string
	Position2  string
	Kind       string
	Name       string
	Definition bool
	ChildNodes []Node
}

func parseRecordDecl(line string) *RecordDecl {
	groups := groupsFromRegex(
		`(?:parent (?P<parent>0x[0-9a-f]+) )?
		(?:prev (?P<prev>0x[0-9a-f]+) )?
		<(?P<position>.*)>
		[ ](?P<position2>[^ ]+ )?
		(?P<kind>struct|union)
		(?P<name>.*)`,
		line,
	)

	definition := false
	name := strings.TrimSpace(groups["name"])
	if name == "definition" {
		name = ""
		definition = true
	}
	if strings.HasSuffix(name, " definition") {
		name = name[0 : len(name)-11]
		definition = true
	}

	return &RecordDecl{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Prev:       groups["prev"],
		Position2:  strings.TrimSpace(groups["position2"]),
		Kind:       groups["kind"],
		Name:       name,
		Definition: definition,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *RecordDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *RecordDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *RecordDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *RecordDecl) Position() Position {
	return n.Pos
}
