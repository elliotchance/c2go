package ast

import "fmt"

type ParenExpr struct {
	Address  string
	Position string
	Type     string
	Children []Node
}

func parseParenExpr(line string) *ParenExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)'",
		line,
	)

	return &ParenExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *ParenExpr) render(ast *Ast) (string, string) {
	a, aType := renderExpression(ast, n.Children[0])
	src := fmt.Sprintf("(%s)", a)
	return src, aType
}
