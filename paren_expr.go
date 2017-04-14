package main

import "fmt"

type ParenExpr struct {
	Address  string
	Position string
	Type     string
	Children []interface{}
}

func parseParenExpr(line string) *ParenExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)'",
		line,
	)

	return &ParenExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Children: []interface{}{},
	}
}

func (n *ParenExpr) Render() []string {
	a := renderExpression(n.Children[0])
	return []string{fmt.Sprintf("(%s)", a[0]), a[1]}
}
