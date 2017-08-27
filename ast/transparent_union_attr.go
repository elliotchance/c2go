package ast

type TransparentUnionAttr struct {
	Addr       Address
	Pos        string
	ChildNodes []Node
}

func parseTransparentUnionAttr(line string) *TransparentUnionAttr {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &TransparentUnionAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        groups["position"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *TransparentUnionAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *TransparentUnionAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *TransparentUnionAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *TransparentUnionAttr) Position() Position {
	return NewPositionFromString(n.Pos)
}
