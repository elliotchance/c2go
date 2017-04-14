package main

import "fmt"

type ArraySubscriptExpr struct {
	Address  string
	Position string
	Type     string
	Kind     string
	Children []interface{}
}

func parseArraySubscriptExpr(line string) *ArraySubscriptExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<kind>.*)",
		line,
	)

	return &ArraySubscriptExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Kind:     groups["kind"],
		Children: []interface{}{},
	}
}

func (n *ArraySubscriptExpr) Render() []string {
	lhs := renderExpression(n.Children[0])
	rhs := renderExpression(n.Children[1])
	newExpression := fmt.Sprintf("%s[%s]", lhs[0], rhs[0])

	newType, err := getDereferenceType(lhs[1])
	if err != nil {
		panic(fmt.Sprintf("Cannot dereference type '%s' for the expression '%s'",
			lhs[1], newExpression))
	}

	return []string{newExpression, newType}
}
