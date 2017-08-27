package ast

type BinaryOperator struct {
	Addr       Address
	Position   string
	Type       string
	Operator   string
	ChildNodes []Node
}

func parseBinaryOperator(line string) *BinaryOperator {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' '(?P<operator>.*?)'",
		line,
	)

	return &BinaryOperator{
		Addr:       ParseAddress(groups["address"]),
		Position:   groups["position"],
		Type:       groups["type"],
		Operator:   groups["operator"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *BinaryOperator) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *BinaryOperator) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *BinaryOperator) Children() []Node {
	return n.ChildNodes
}
