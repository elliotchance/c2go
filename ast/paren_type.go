package ast

type ParenType struct {
	Addr     Address
	Type     string
	Sugar    bool
	Children []Node
}

func parseParenType(line string) *ParenType {
	groups := groupsFromRegex(`'(?P<type>.*?)' sugar`, line)

	return &ParenType{
		Addr:     ParseAddress(groups["address"]),
		Type:     groups["type"],
		Sugar:    true,
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ParenType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
