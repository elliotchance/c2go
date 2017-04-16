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
	expression, _ := renderExpression(ast, children[0])
	index, _ := renderExpression(ast, children[1])
	src := fmt.Sprintf("%s[%s]", expression, index)
	return src, "unknown1"
}
