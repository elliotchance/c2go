package ast

import "fmt"

type ArraySubscriptExpr struct {
	Address  string
	Position string
	Type     string
	Kind     string
	Children []Node
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
		Children: []Node{},
	}
}

func (n *ArraySubscriptExpr) render(ast *Ast) (string, string) {
	children := n.Children
	expression, expressionType := renderExpression(ast, children[0])
	index, _ := renderExpression(ast, children[1])
	src := fmt.Sprintf("%s[%s]", expression, index)

	newType, err := getDereferenceType(expressionType)
	if err != nil {
		panic(fmt.Sprintf("Cannot dereference type '%s' for the expression '%s'",
			expressionType, expression))
	}

	return src, newType
}

func (n *ArraySubscriptExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
