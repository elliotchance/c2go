package ast

type PackedAttr struct {
	Addr     Address
	Position string
	Children []Node
}

func parsePackedAttr(line string) *PackedAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &PackedAttr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *PackedAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
