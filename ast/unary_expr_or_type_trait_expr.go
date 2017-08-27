package ast

type UnaryExprOrTypeTraitExpr struct {
	Addr     Address
	Position string
	Type1    string
	Function string
	Type2    string
	Children []Node
}

func parseUnaryExprOrTypeTraitExpr(line string) *UnaryExprOrTypeTraitExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type1>.+?)'
		 (?P<function>[^ ]+)
		(?P<type2> '.+?')?`,
		line,
	)

	return &UnaryExprOrTypeTraitExpr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Type1:    groups["type1"],
		Function: groups["function"],
		Type2:    removeQuotes(groups["type2"]),
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *UnaryExprOrTypeTraitExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
