package main

import "fmt"

type BinaryOperator struct {
	Address  string
	Position string
	Type     string
	Operator string
	Children []interface{}
}

func parseBinaryOperator(line string) *BinaryOperator {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' '(?P<operator>.*?)'",
		line,
	)

	return &BinaryOperator{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Operator: groups["operator"],
		Children: []interface{}{},
	}
}

func (n *BinaryOperator) Render() []string {
	operator := n.Operator

	left := renderExpression(n.Children[0])
	right := renderExpression(n.Children[1])

	return_type := "bool"
	if inStrings(operator, []string{"|", "&", "+", "-", "*", "/"}) {
		// TODO: The left and right type might be different
		return_type = left[1]
	}

	if operator == "&&" {
		left[0] = cast(left[0], left[1], return_type)
		right[0] = cast(right[0], right[1], return_type)
	}

	if (operator == "!=" || operator == "==") && right[0] == "(0)" {
		right[0] = "nil"
	}

	return []string{fmt.Sprintf("%s %s %s", left[0], operator, right[0]), return_type}
}
