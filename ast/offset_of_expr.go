package ast

type OffsetOfExpr struct {
	Address  string
	Position string
	Type     string
	Children []Node
}

func parseOffsetOfExpr(line string) *OffsetOfExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*)'",
		line,
	)

	return &OffsetOfExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *OffsetOfExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
