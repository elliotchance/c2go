package ast

import (
	"fmt"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
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

	return_type := "bool"
	if util.InStrings(operator, []string{"|", "&", "+", "-", "*", "/", "=", "<<", ">>"}) {
		// TODO: The left and right type might be different
		return_type = leftType
	}

	if operator == "&&" {
		left = types.Cast(program, left, leftType, return_type)
		right = types.Cast(program, right, rightType, return_type)
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
