package ast

type VAArgExpr struct {
	Address  string
	Position string
	Type     string
	Children []Node
}

func parseVAArgExpr(line string) *VAArgExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*)'",
		line,
	)

	return &VAArgExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *VAArgExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
