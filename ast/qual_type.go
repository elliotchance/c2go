package ast

type QualType struct {
	Addr     Address
	Type     string
	Kind     string
	Children []Node
}

func parseQualType(line string) *QualType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<kind>.*)",
		line,
	)

	return &QualType{
		Addr:     ParseAddress(groups["address"]),
		Type:     groups["type"],
		Kind:     groups["kind"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *QualType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
