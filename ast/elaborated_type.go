package ast

type ElaboratedType struct {
	Address  string
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
		Address:  groups["address"],
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
