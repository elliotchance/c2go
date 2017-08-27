package ast

type NoInlineAttr struct {
	Addr       Address
	Pos        string
	ChildNodes []Node
}

func parseNoInlineAttr(line string) *NoInlineAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &NoInlineAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        groups["position"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *NoInlineAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *NoInlineAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *NoInlineAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *NoInlineAttr) Position() Position {
	return NewPositionFromString(n.Pos)
}
