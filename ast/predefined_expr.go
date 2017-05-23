package ast

type PredefinedExpr struct {
	Address  string
	Position string
	Type     string
	Name     string
	Lvalue   bool
	Children []Node
}

func parsePredefinedExpr(line string) *PredefinedExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*)' lvalue (?P<name>.*)",
		line,
	)

	return &PredefinedExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Name:     groups["name"],
		Lvalue:   true,
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *PredefinedExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
