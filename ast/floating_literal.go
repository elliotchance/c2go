package ast

type FloatingLiteral struct {
	Addr     Address
	Position string
	Type     string
	Value    float64
	Children []Node
}

func parseFloatingLiteral(line string) *FloatingLiteral {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<value>.+)",
		line,
	)

	return &FloatingLiteral{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Type:     groups["type"],
		Value:    atof(groups["value"]),
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FloatingLiteral) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *FloatingLiteral) Address() Address {
	return n.Addr
}
