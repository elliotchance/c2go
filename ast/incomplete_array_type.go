package ast

type IncompleteArrayType struct {
	Address  string
	Type     string
	Children []Node
}

func parseIncompleteArrayType(line string) *IncompleteArrayType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' ",
		line,
	)

	return &IncompleteArrayType{
		Address:  groups["address"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *IncompleteArrayType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
