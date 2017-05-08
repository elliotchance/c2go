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

func (n *ElaboratedType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
