package ast

type ArraySubscriptExpr struct {
	Addr       Address
	Pos        string
	Type       string
	Kind       string
	ChildNodes []Node
}

func parseArraySubscriptExpr(line string) *ArraySubscriptExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<kind>.*)",
		line,
	)

	return &ArraySubscriptExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        groups["position"],
		Type:       groups["type"],
		Kind:       groups["kind"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ArraySubscriptExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ArraySubscriptExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ArraySubscriptExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ArraySubscriptExpr) Position() Position {
	return NewPositionFromString(n.Pos)
}
