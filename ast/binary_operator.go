package ast

type BinaryOperator struct {
	Addr     Address
	Position string
	Type     string
	Operator string
	Children []Node
}

func parseBinaryOperator(line string) *BinaryOperator {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' '(?P<operator>.*?)'",
		line,
	)

	return &BinaryOperator{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Type:     groups["type"],
		Operator: groups["operator"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *BinaryOperator) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
