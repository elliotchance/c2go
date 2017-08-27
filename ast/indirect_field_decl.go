package ast

import "strings"

type IndirectFieldDecl struct {
	Addr      Address
	Position  string
	Position2 string
	Implicit  bool
	Name      string
	Type      string
	Children  []Node
}

func parseIndirectFieldDecl(line string) *IndirectFieldDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<position2> [^ ]+:[\d:]+)?
		(?P<implicit> implicit)?
		 (?P<name>\w+)
		 '(?P<type>.+)'`,
		line,
	)

	return &IndirectFieldDecl{
		Addr:      ParseAddress(groups["address"]),
		Position:  groups["position"],
		Position2: strings.TrimSpace(groups["position2"]),
		Implicit:  len(groups["implicit"]) > 0,
		Name:      groups["name"],
		Type:      groups["type"],
		Children:  []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *IndirectFieldDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
