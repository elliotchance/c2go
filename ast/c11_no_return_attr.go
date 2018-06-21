package ast

// C11NoReturnAttr is a type of attribute that is optionally attached to a function
// with return type void.
type C11NoReturnAttr struct {
	Addr       Address
	Pos        Position
	ChildNodes []Node
}

func parseC11NoReturnAttr(line string) *C11NoReturnAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &C11NoReturnAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *C11NoReturnAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *C11NoReturnAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *C11NoReturnAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *C11NoReturnAttr) Position() Position {
	return n.Pos
}
