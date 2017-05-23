package ast

type RestrictAttr struct {
	Address  string
	Position string
	Name     string
	Children []Node
}

func parseRestrictAttr(line string) *RestrictAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> (?P<name>.+)",
		line,
	)

	return &RestrictAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Name:     groups["name"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *RestrictAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
