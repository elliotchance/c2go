package ast

type ImplicitCastExpr struct {
	Address  string
	Position string
	Type     string
	Kind     string
	Children []Node
}

func parseImplicitCastExpr(line string) *ImplicitCastExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*)' <(?P<kind>.*)>",
		line,
	)

	return &ImplicitCastExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Kind:     groups["kind"],
		Children: []Node{},
	}
}

func (n *ImplicitCastExpr) render(ast *Ast) (string, string) {
	return renderExpression(ast, n.Children[0])
}
