package ast

type CStyleCastExpr struct {
	Addr       Address
	Position   string
	Type       string
	Kind       string
	ChildNodes []Node
}

func parseCStyleCastExpr(line string) *CStyleCastExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' <(?P<kind>.*)>",
		line,
	)

	return &CStyleCastExpr{
		Addr:       ParseAddress(groups["address"]),
		Position:   groups["position"],
		Type:       groups["type"],
		Kind:       groups["kind"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CStyleCastExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CStyleCastExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CStyleCastExpr) Children() []Node {
	return n.ChildNodes
}
