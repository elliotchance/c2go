package ast

import "fmt"

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

func (n *ConditionalOperator) render(ast *Ast) (string, string) {
	a, _ := renderExpression(ast, n.Children[0])
	b, _ := renderExpression(ast, n.Children[1])
	c, _ := renderExpression(ast, n.Children[2])

	ast.addImport("github.com/elliotchance/c2go/noarch")
	src := fmt.Sprintf("noarch.Ternary(%s, func () interface{} { return %s }, func () interface{} { return %s })", a, b, c)
	return src, n.Type
}

func (n *ConditionalOperator) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
