package ast

// TypedefType is typedef type
type TypedefType struct {
	Addr       Address
	Type       string
	Tags       string
	ChildNodes []Node
}

func parseTypedefType(line string) *TypedefType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<tags>.+)",
		line,
	)

	return &TypedefType{
		Addr:       ParseAddress(groups["address"]),
		Type:       groups["type"],
		Tags:       groups["tags"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *TypedefType) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *TypedefType) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *TypedefType) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *TypedefType) Position() Position {
	return Position{}
}
