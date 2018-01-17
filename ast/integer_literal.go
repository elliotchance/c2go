package ast

// IntegerLiteral is type of integer literal
type IntegerLiteral struct {
	Addr       Address
	Pos        Position
	Type       string
	Value      string
	ChildNodes []Node
}

func parseIntegerLiteral(line string) *IntegerLiteral {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<value>\\d+)",
		line,
	)

	return &IntegerLiteral{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Value:      groups["value"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *IntegerLiteral) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *IntegerLiteral) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *IntegerLiteral) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *IntegerLiteral) Position() Position {
	return n.Pos
}
