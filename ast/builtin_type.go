package ast

type BuiltinType struct {
	Address  string
	Type     string
	Children []Node
}

func parseBuiltinType(line string) *BuiltinType {
	groups := groupsFromRegex(
		"'(?P<type>.*?)'",
		line,
	)

	return &BuiltinType{
		Address:  groups["address"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *BuiltinType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
