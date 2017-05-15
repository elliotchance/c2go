package ast

type AlignedAttr struct {
	Address  string
	Position string
	Children []Node
}

func parseAlignedAttr(line string) *AlignedAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> aligned",
		line,
	)

	return &AlignedAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *AlignedAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
