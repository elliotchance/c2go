package ast

// AlwaysInlineAttr is a type of attribute that is optionally attached to a
// variable or struct field definition.
type AlwaysInlineAttr struct {
	Addr     Address
	Position string
	Children []Node
}

func parseAlwaysInlineAttr(line string) *AlwaysInlineAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> always_inline",
		line,
	)

	return &AlwaysInlineAttr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *AlwaysInlineAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
