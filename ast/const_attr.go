package ast

// ConstAttr is a type of attribute that is optionally attached to a variable
// or struct field definition.
type ConstAttr struct {
	Addr       Address
	Pos        Position
	Tags       string
	ChildNodes []Node
}

func parseConstAttr(line string) *ConstAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>(?P<tags>.*)",
		line,
	)

	return &ConstAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Tags:       groups["tags"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ConstAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ConstAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ConstAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ConstAttr) Position() Position {
	return n.Pos
}
