package ast

type InitListExpr struct {
	Address  string
	Position string
	Type     string
	Children []Node
}

func parseInitListExpr(line string) *InitListExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*)'",
		line,
	)

	return &InitListExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *InitListExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
