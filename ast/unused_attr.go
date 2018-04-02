package ast

// UnusedAttr is a type of attribute that is optionally attached to a variable
// or struct field definition.
type UnusedAttr struct {
	Addr       Address
	Pos        Position
	ChildNodes []Node
	IsUnused   bool
}

func parseUnusedAttr(line string) *UnusedAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>(?P<unused> unused)?",
		line,
	)

	return &UnusedAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		ChildNodes: []Node{},
		IsUnused:   len(groups["unused"]) > 0,
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *UnusedAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *UnusedAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *UnusedAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *UnusedAttr) Position() Position {
	return n.Pos
}
