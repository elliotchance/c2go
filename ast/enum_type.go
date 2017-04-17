package ast

type EnumType struct {
	Address  string
	Name     string
	Children []Node
}

func parseEnumType(line string) *EnumType {
	groups := groupsFromRegex(
		"'(?P<name>.*)'",
		line,
	)

	return &EnumType{
		Address:  groups["address"],
		Name:     groups["name"],
		Children: []Node{},
	}
}

func (n *EnumType) render(ast *Ast) (string, string) {
	return "", ""
}

func (n *EnumType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
