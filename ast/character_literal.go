package ast

import (
	"github.com/elliotchance/c2go/util"
)

type CharacterLiteral struct {
	Addr     Address
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
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Type:     groups["type"],
		Value:    util.Atoi(groups["value"]),
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CharacterLiteral) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CharacterLiteral) Address() Address {
	return n.Addr
}
