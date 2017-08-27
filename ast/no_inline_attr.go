package ast

type NoInlineAttr struct {
	Addr     Address
	Position string
	Children []Node
}

func parseNoInlineAttr(line string) *NoInlineAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &NoInlineAttr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *NoInlineAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *NoInlineAttr) Address() Address {
	return n.Addr
}
