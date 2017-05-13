package ast

import (
	"fmt"
	"strconv"
)

type StringLiteral struct {
	Address  string
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
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Value:    s,
		Lvalue:   true,
		Children: []Node{},
	}
}

func (n *StringLiteral) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
