package ast

// AlwaysInlineAttr is a type of attribute that is optionally attached to a
// variable or struct field definition.
type AlwaysInlineAttr struct {
	Addr       Address
	Pos        Position
	ChildNodes []Node
}

func parseAlwaysInlineAttr(line string) *AlwaysInlineAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> always_inline",
		line,
	)

	return &AlwaysInlineAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *AlwaysInlineAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *AlwaysInlineAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *AlwaysInlineAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *AlwaysInlineAttr) Position() Position {
	return n.Pos
}
