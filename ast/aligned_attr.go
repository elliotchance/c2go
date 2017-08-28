package ast

// AlignedAttr is a type of attribute that is optionally attached to a variable
// or struct field definition.
type AlignedAttr struct {
	Addr       Address
	Pos        Position
	IsAligned  bool
	ChildNodes []Node
}

func parseAlignedAttr(line string) *AlignedAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>(?P<aligned> aligned)?",
		line,
	)

	return &AlignedAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		IsAligned:  len(groups["aligned"]) > 0,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *AlignedAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *AlignedAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *AlignedAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *AlignedAttr) Position() Position {
	return n.Pos
}
