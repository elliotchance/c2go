package ast

type UnusedAttr struct {
	Addr       Address
	Pos        string
	ChildNodes []Node
}

func parseUnusedAttr(line string) *UnusedAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> unused",
		line,
	)

	return &UnusedAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        groups["position"],
		ChildNodes: []Node{},
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
	return NewPositionFromString(n.Pos)
}
