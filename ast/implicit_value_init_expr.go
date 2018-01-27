package ast

// ImplicitValueInitExpr is expression
type ImplicitValueInitExpr struct {
	Addr       Address
	Pos        Position
	Type1      string
	Type2      string
	ChildNodes []Node
}

func parseImplicitValueInitExpr(line string) *ImplicitValueInitExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type1>.*?)'(:'(?P<type2>.*)')?",
		line,
	)

	return &ImplicitValueInitExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type1:      groups["type1"],
		Type2:      groups["type2"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ImplicitValueInitExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ImplicitValueInitExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ImplicitValueInitExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ImplicitValueInitExpr) Position() Position {
	return n.Pos
}
