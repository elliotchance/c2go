package ast

import (
	"fmt"

	"github.com/elliotchance/c2go/program"
)

type ConditionalOperator struct {
	Address  string
	Position string
	Type     string
	Children []Node
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
		Children: []Node{},
	}
}

func (n *ConditionalOperator) render(program *program.Program) (string, string) {
	a, _ := renderExpression(program, n.Children[0])
	b, _ := renderExpression(program, n.Children[1])
	c, _ := renderExpression(program, n.Children[2])

	program.AddImport("github.com/elliotchance/c2go/noarch")
	src := fmt.Sprintf("noarch.Ternary(%s, func () interface{} { return %s }, func () interface{} { return %s })", a, b, c)
	return src, n.Type
}

func (n *ConditionalOperator) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
