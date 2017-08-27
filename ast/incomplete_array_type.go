package ast

type IncompleteArrayType struct {
	Addr     Address
	Type     string
	Children []Node
}

func parseIncompleteArrayType(line string) *IncompleteArrayType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' ",
		line,
	)

	return &IncompleteArrayType{
		Addr:     ParseAddress(groups["address"]),
		Type:     groups["type"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *IncompleteArrayType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
