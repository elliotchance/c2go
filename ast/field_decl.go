package ast

import (
	"strings"
)

type FieldDecl struct {
	Address    string
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
		Address:    groups["address"],
		Position:   groups["position"],
		Position2:  strings.TrimSpace(groups["position2"]),
		Name:       strings.TrimSpace(groups["name"]),
		Type:       groups["type"],
		Referenced: len(groups["referenced"]) > 0,
		Children:   []Node{},
	}
}

func (n *FieldDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
