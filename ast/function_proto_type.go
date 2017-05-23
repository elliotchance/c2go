package ast

type FunctionProtoType struct {
	Address  string
	Type     string
	Kind     string
	Children []Node
}

func parseFunctionProtoType(line string) *FunctionProtoType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<kind>.*)",
		line,
	)

	return &FunctionProtoType{
		Address:  groups["address"],
		Type:     groups["type"],
		Kind:     groups["kind"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FunctionProtoType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
