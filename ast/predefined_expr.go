package ast

// PredefinedExpr is expression.
type PredefinedExpr struct {
	Addr       Address
	Pos        Position
	Type       string
	Name       string
	Lvalue     bool
	ChildNodes []Node
}

func parsePredefinedExpr(line string) *PredefinedExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*)' lvalue (?P<name>.*)",
		line,
	)

	return &PredefinedExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Name:       groups["name"],
		Lvalue:     true,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *PredefinedExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *PredefinedExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *PredefinedExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *PredefinedExpr) Position() Position {
	return n.Pos
}
