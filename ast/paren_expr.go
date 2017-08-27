package ast

type ParenExpr struct {
	Addr       Address
	Pos        string
	Type       string
	ChildNodes []Node
}

func parseParenExpr(line string) *ParenExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)'",
		line,
	)

	return &ParenExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        groups["position"],
		Type:       groups["type"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ParenExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ParenExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ParenExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ParenExpr) Position() Position {
	return NewPositionFromString(n.Pos)
}
