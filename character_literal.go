package main

import "fmt"

type CharacterLiteral struct {
	Address  string
	Position string
	Type     string
	Value    int
	Children []interface{}
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
		Value:    atoi(groups["value"]),
		Children: []interface{}{},
	}
}

func (n *CharacterLiteral) Render() []string {
	var s string

	switch n.Value {
	case '\n':
		s = "'\\n'"
	default:
		s = fmt.Sprintf("'%c'", n.Value)
	}

	return []string{s, n.Type}
}
