package ast

type InitListExpr struct {
	Addr     Address
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
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *InitListExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
