package ast

type UnusedAttr struct {
	Addr     Address
	Position string
	Children []Node
}

func parseUnusedAttr(line string) *UnusedAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> unused",
		line,
	)

	return &UnusedAttr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *UnusedAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *UnusedAttr) Address() Address {
	return n.Addr
}
