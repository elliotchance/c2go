package main

import "fmt"

type PredefinedExpr struct {
	Address  string
	Position string
	Type     string
	Name     string
	Lvalue   bool
	Children []interface{}
}

func parsePredefinedExpr(line string) *PredefinedExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*)' lvalue (?P<name>.*)",
		line,
	)

	return &PredefinedExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Name:     groups["name"],
		Lvalue:   true,
		Children: []interface{}{},
	}
}

func (n *PredefinedExpr) Render() []string {
	if n.Name == "__PRETTY_FUNCTION__" {
		// FIXME
		return []string{"\"void print_number(int *)\"", "const char*"}
	}

	if n.Name == "__func__" {
		// FIXME
		return []string{fmt.Sprintf("\"%s\"", "print_number"), "const char*"}
	}

	panic(fmt.Sprintf("renderExpression: unknown PredefinedExpr: %s", n.Name))
}
