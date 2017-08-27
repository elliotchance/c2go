package ast

type TypedefType struct {
	Addr     Address
	Type     string
	Tags     string
	Children []Node
}

func parseTypedefType(line string) *TypedefType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<tags>.+)",
		line,
	)

	return &TypedefType{
		Addr:     ParseAddress(groups["address"]),
		Type:     groups["type"],
		Tags:     groups["tags"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *TypedefType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
