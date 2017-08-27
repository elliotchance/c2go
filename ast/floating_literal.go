package ast

type FloatingLiteral struct {
	Addr       Address
	Position   string
	Type       string
	Value      float64
	ChildNodes []Node
}

func parseFloatingLiteral(line string) *FloatingLiteral {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<value>.+)",
		line,
	)

	return &FloatingLiteral{
		Addr:       ParseAddress(groups["address"]),
		Position:   groups["position"],
		Type:       groups["type"],
		Value:      atof(groups["value"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FloatingLiteral) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *FloatingLiteral) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *FloatingLiteral) Children() []Node {
	return n.ChildNodes
}
