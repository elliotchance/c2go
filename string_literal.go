package main

import (
	"fmt"
	"strings"
)

type StringLiteral struct {
	Address  string
	Position string
	Type     string
	Value    string
	Lvalue   bool
	Children []interface{}
}

func parseStringLiteral(line string) *StringLiteral {
	groups := groupsFromRegex(
		`<(?P<position>.*)> '(?P<type>.*)' lvalue "(?P<value>.*)"`,
		line,
	)

	return &StringLiteral{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Value:    unescapeString(groups["value"]),
		Lvalue:   true,
		Children: []interface{}{},
	}
}

func (n *StringLiteral) Render() []string {
	return []string{
		fmt.Sprintf("\"%s\"", strings.Replace(n.Value, "\n", "\\n", -1)),
		"const char *",
	}
}
