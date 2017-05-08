package ast

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

func (n *ParenExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
