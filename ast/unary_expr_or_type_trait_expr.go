package ast

type UnaryExprOrTypeTraitExpr struct {
	Addr       Address
	Position   string
	Type1      string
	Function   string
	Type2      string
	ChildNodes []Node
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
		Addr:       ParseAddress(groups["address"]),
		Position:   groups["position"],
		Type1:      groups["type1"],
		Function:   groups["function"],
		Type2:      removeQuotes(groups["type2"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *UnaryExprOrTypeTraitExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *UnaryExprOrTypeTraitExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *UnaryExprOrTypeTraitExpr) Children() []Node {
	return n.ChildNodes
}
