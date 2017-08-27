package ast

type EnumType struct {
	Addr     Address
	Name     string
	Children []Node
}

func parseEnumType(line string) *EnumType {
	groups := groupsFromRegex(
		"'(?P<name>.*)'",
		line,
	)

	return &EnumType{
		Addr:     ParseAddress(groups["address"]),
		Name:     groups["name"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *EnumType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *EnumType) Address() Address {
	return n.Addr
}
