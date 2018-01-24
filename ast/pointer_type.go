package ast

// PointerType is pointer type
type PointerType struct {
	Addr       Address
	Type       string
	ChildNodes []Node
}

func parsePointerType(line string) *PointerType {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return &PointerType{
		Addr:       ParseAddress(groups["address"]),
		Type:       groups["type"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *PointerType) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *PointerType) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *PointerType) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *PointerType) Position() Position {
	return Position{}
}
