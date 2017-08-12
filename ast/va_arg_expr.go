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

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *VAArgExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
