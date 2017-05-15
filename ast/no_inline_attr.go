package ast

type NoInlineAttr struct {
	Address  string
	Position string
	Children []Node
}

func parseNoInlineAttr(line string) *NoInlineAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &NoInlineAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *NoInlineAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
