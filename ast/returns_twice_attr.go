package ast

type ReturnsTwiceAttr struct {
	Addr     Address
	Position string
	Children []Node
}

func parseReturnsTwiceAttr(line string) *ReturnsTwiceAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ReturnsTwiceAttr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ReturnsTwiceAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
