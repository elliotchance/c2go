package ast

type ParenExpr struct {
	Addr     Address
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
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ParenExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ParenExpr) Address() Address {
	return n.Addr
}
