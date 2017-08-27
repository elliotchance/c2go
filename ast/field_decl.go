package ast

import (
	"strings"
)

type FieldDecl struct {
	Addr       Address
	Position   string
	Position2  string
	Name       string
	Type       string
	Referenced bool
	Children   []Node
}

func parseFieldDecl(line string) *FieldDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<position2> col:\d+| line:\d+:\d+)?
		(?P<referenced> referenced)?
		(?P<name> \w+?)?
		 '(?P<type>.+?)'`,
		line,
	)

	return &FieldDecl{
		Addr:       ParseAddress(groups["address"]),
		Position:   groups["position"],
		Position2:  strings.TrimSpace(groups["position2"]),
		Name:       strings.TrimSpace(groups["name"]),
		Type:       groups["type"],
		Referenced: len(groups["referenced"]) > 0,
		Children:   []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FieldDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
