package ast

// AlignedAttr is a type of attribute that is optionally attached to a variable
// or struct field definition.
type AlignedAttr struct {
	Addr      Address
	Position  string
	IsAligned bool
	Children  []Node
}

func parseAlignedAttr(line string) *AlignedAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>(?P<aligned> aligned)?",
		line,
	)

	return &AlignedAttr{
		Addr:      ParseAddress(groups["address"]),
		Position:  groups["position"],
		IsAligned: len(groups["aligned"]) > 0,
		Children:  []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *AlignedAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
