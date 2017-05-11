package ast

type IntegerLiteral struct {
	Address  string
	Position string
	Type     string
	Value    string
	Children []Node
}

func parseIntegerLiteral(line string) *IntegerLiteral {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<value>\\d+)",
		line,
	)

	return &IntegerLiteral{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Value:    groups["value"],
		Children: []Node{},
	}
}

func (n *IntegerLiteral) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
