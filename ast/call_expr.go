package ast

type CallExpr struct {
	Address  string
	Position string
	Type     string
	Children []Node
}

func parseCallExpr(line string) *CallExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)'",
		line,
	)

	return &CallExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *CallExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
