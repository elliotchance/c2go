package ast

type UnusedAttr struct {
	Address  string
	Position string
	Children []Node
}

func parseUnusedAttr(line string) *UnusedAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> unused",
		line,
	)

	return &UnusedAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *UnusedAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
