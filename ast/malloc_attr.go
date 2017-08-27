package ast

type MallocAttr struct {
	Addr       Address
	Position   string
	ChildNodes []Node
}

func parseMallocAttr(line string) *MallocAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &MallocAttr{
		Addr:       ParseAddress(groups["address"]),
		Position:   groups["position"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *MallocAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *MallocAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *MallocAttr) Children() []Node {
	return n.ChildNodes
}
