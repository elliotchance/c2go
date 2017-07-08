package ast

type MallocAttr struct {
	Address  string
	Position string
	Children []Node
}

func parseMallocAttr(line string) *MallocAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &MallocAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *MallocAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
