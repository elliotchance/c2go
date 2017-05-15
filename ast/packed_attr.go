package ast

type PackedAttr struct {
	Address  string
	Position string
	Children []Node
}

func parsePackedAttr(line string) *PackedAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &PackedAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *PackedAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
