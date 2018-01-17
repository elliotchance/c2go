package ast

// QualType is qual type
type QualType struct {
	Addr       Address
	Type       string
	Kind       string
	ChildNodes []Node
}

func parseQualType(line string) *QualType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<kind>.*)",
		line,
	)

	return &QualType{
		Addr:       ParseAddress(groups["address"]),
		Type:       groups["type"],
		Kind:       groups["kind"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *QualType) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *QualType) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *QualType) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *QualType) Position() Position {
	return Position{}
}
