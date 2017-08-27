package ast

type ConstAttr struct {
	Addr     Address
	Position string
	Tags     string
	Children []Node
}

func parseConstAttr(line string) *ConstAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>(?P<tags>.*)",
		line,
	)

	return &ConstAttr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Tags:     groups["tags"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ConstAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
