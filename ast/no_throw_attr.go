package ast

type NoThrowAttr struct {
	Address  string
	Position string
	Children []Node
}

func parseNoThrowAttr(line string) *NoThrowAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &NoThrowAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *NoThrowAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
