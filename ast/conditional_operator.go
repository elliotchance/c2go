package ast

type ConditionalOperator struct {
	Address  string
	Position string
	Type     string
	Children []Node
}

func parseConditionalOperator(line string) *ConditionalOperator {
	groups := groupsFromRegex(
		`<(?P<position>.*)> '(?P<type>.*?)'`,
		line,
	)

	return &ConditionalOperator{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *ConditionalOperator) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
