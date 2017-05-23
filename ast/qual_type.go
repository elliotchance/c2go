package ast

type QualType struct {
	Address  string
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
		Address:  groups["address"],
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
