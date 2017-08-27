package ast

type CallExpr struct {
	Addr     Address
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
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CallExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
