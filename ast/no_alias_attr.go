package ast

// NoAliasAttr is a type of attribute that is optionally attached to a function declaration.
type NoAliasAttr struct {
	Addr       Address
	Pos        Position
	ChildNodes []Node
}

func parseNoAliasAttr(line string) *NoAliasAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &NoAliasAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *NoAliasAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *NoAliasAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *NoAliasAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *NoAliasAttr) Position() Position {
	return n.Pos
}
