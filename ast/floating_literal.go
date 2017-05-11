package ast

type FloatingLiteral struct {
	Address  string
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
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Value:    atof(groups["value"]),
		Children: []Node{},
	}
}

func (n *FloatingLiteral) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
