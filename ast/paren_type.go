package ast

type ParenType struct {
	Address  string
	Type     string
	Sugar    bool
	Children []interface{}
}

func parseParenType(line string) *ParenType {
	groups := groupsFromRegex(`'(?P<type>.*?)' sugar`, line)

	return &ParenType{
		Address:  groups["address"],
		Type:     groups["type"],
		Sugar:    true,
		Children: []interface{}{},
	}
}

func (n *ParenType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
