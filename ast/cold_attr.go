package ast

type ColdAttr struct {
	Addr       Address
	Pos        Position
	ChildNodes []Node
}

func parseColdAttr(line string) *ColdAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>`,
		line,
	)

	return &ColdAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ColdAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ColdAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ColdAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ColdAttr) Position() Position {
	return n.Pos
}
