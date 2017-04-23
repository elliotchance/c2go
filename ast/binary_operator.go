package ast

import (
	"fmt"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
)

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

func (n *BinaryOperator) render(program *program.Program) (string, string) {
	operator := n.Operator
	left, leftType := renderExpression(program, n.Children[0])
	right, rightType := renderExpression(program, n.Children[1])
	returnType := types.ResolveTypeForBinaryOperator(program, operator, leftType, rightType)

	if operator == "&&" {
		left = types.Cast(program, left, leftType, "bool")
		right = types.Cast(program, right, rightType, "bool")

		src := fmt.Sprintf("%s %s %s", left, operator, right)
		return types.Cast(program, src, "bool", returnType), returnType
	}

	if (operator == "!=" || operator == "==") && right == "(0)" {
		right = "nil"
	}

	if operator == "=" {
		right = types.Cast(program, right, rightType, returnType)
	}

	src := fmt.Sprintf("%s %s %s", left, operator, right)
	return src, returnType
}

func (n *BinaryOperator) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
