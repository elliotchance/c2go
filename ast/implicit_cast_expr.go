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

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ImplicitCastExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
