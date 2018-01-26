package ast

// ParenType is paren type
type ParenType struct {
	Addr       Address
	Type       string
	Sugar      bool
	ChildNodes []Node
}

func parseParenType(line string) *ParenType {
	groups := groupsFromRegex(`'(?P<type>.*?)' sugar`, line)

	return &ParenType{
		Addr:       ParseAddress(groups["address"]),
		Type:       groups["type"],
		Sugar:      true,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ParenType) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ParenType) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ParenType) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ParenType) Position() Position {
	return Position{}
}
