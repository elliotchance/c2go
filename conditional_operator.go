package main

import "fmt"

type ConditionalOperator struct {
	Address  string
	Position string
	Type     string
	Children []interface{}
}

func parseConditionalOperator(line string) *ConditionalOperator {
	groups := groupsFromRegex(
		`<(?P<position>.*)> '(?P<type>.*?)'`,
		line,
	)

	return &ConditionalOperator{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Children: []interface{}{},
	}
}

func (n *ConditionalOperator) Render() []string {
	a := renderExpression(n.Children[0])[0]
	b := renderExpression(n.Children[1])[0]
	c := renderExpression(n.Children[2])[0]

	addImport("github.com/elliotchance/c2go/noarch")
	return []string{
		fmt.Sprintf("noarch.Ternary(%s, func () interface{} { return %s }, func () interface{} { return %s })", a, b, c),
		n.Type,
	}
}
