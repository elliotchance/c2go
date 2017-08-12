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

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *NoThrowAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
