package ast

import (
	"github.com/elliotchance/c2go/util"
)

type CharacterLiteral struct {
	Address  string
	Position string
	Type     string
	Value    int
	Children []Node
}

func parseCharacterLiteral(line string) *CharacterLiteral {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<value>\\d+)",
		line,
	)

	return &CharacterLiteral{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Value:    util.Atoi(groups["value"]),
		Children: []Node{},
	}
}

func (n *CharacterLiteral) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
