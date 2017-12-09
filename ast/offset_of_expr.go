package ast

// OffsetOfExpr is expression.
type OffsetOfExpr struct {
	Addr       Address
	Pos        Position
	Type       string
	ChildNodes []Node
}

func parseOffsetOfExpr(line string) *OffsetOfExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*)'",
		line,
	)

	return &OffsetOfExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *OffsetOfExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *OffsetOfExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *OffsetOfExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *OffsetOfExpr) Position() Position {
	return n.Pos
}
