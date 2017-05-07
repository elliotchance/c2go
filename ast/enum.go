package ast

type Enum struct {
	Address  string
	Name     string
	Children []Node
}

func parseEnum(line string) *Enum {
	groups := groupsFromRegex(
		"'(?P<name>.*)'",
		line,
	)

	return &Enum{
		Address:  groups["address"],
		Name:     groups["name"],
		Children: []Node{},
	}
}

func (n *Enum) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
