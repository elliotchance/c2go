package ast

import "fmt"

type BinaryOperator struct {
	Address  string
	Position string
	Type     string
	Operator string
	Children []Node
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
		Children: []Node{},
	}
}

func (n *BinaryOperator) render(ast *Ast) (string, string) {
	operator := n.Operator

	left, leftType := renderExpression(ast, n.Children[0])
	right, rightType := renderExpression(ast, n.Children[1])

	return_type := "bool"
	if inStrings(operator, []string{"|", "&", "+", "-", "*", "/", "="}) {
		// TODO: The left and right type might be different
		return_type = leftType
	}

	if operator == "&&" {
		left = cast(ast, left, leftType, return_type)
		right = cast(ast, right, rightType, return_type)
	}

	if (operator == "!=" || operator == "==") && right == "(0)" {
		right = "nil"
	}

	src := fmt.Sprintf("%s %s %s", left, operator, right)
	return src, return_type
}

func (n *BinaryOperator) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
