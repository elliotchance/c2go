package ast

type ElaboratedType struct {
	Addr     Address
	Type     string
	Tags     string
	Children []Node
}

func parseElaboratedType(line string) *ElaboratedType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<tags>.+)",
		line,
	)

	return &ElaboratedType{
		Addr:     ParseAddress(groups["address"]),
		Type:     groups["type"],
		Tags:     groups["tags"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ElaboratedType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ElaboratedType) Address() Address {
	return n.Addr
}
