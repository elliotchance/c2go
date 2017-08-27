package ast

import (
	"strings"
)

type EnumDecl struct {
	Addr      Address
	Position  string
	Position2 string
	Name      string
	Children  []Node
}

func parseEnumDecl(line string) *EnumDecl {
	groups := groupsFromRegex(
		"<(?P<position>.*)>(?P<position2> 0x[^ ]+)?(?P<name>.*)",
		line,
	)

	return &EnumDecl{
		Addr:      ParseAddress(groups["address"]),
		Position:  groups["position"],
		Position2: groups["position2"],
		Name:      strings.TrimSpace(groups["name"]),
		Children:  []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *EnumDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
