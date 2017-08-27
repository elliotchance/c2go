package ast

import (
	"fmt"
	"strconv"
)

type StringLiteral struct {
	Addr     Address
	Position string
	Type     string
	Value    string
	Lvalue   bool
	Children []Node
}

func parseStringLiteral(line string) *StringLiteral {
	groups := groupsFromRegex(
		`<(?P<position>.*)> '(?P<type>.*)' lvalue (?P<value>".*")`,
		line,
	)

	s, err := strconv.Unquote(groups["value"])
	if err != nil {
		panic(fmt.Sprintf("Unable to unquote %s\n", groups["value"]))
	}

	return &StringLiteral{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Type:     groups["type"],
		Value:    s,
		Lvalue:   true,
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *StringLiteral) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
