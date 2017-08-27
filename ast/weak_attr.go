package ast

// WeakAttr for the WeakAttr node
type WeakAttr struct {
	Addr     Address
	Position string
	Children []Node
}

func parseWeakAttr(line string) *WeakAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>`,
		line,
	)

	return &WeakAttr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild method to implements Node interface
func (a *WeakAttr) AddChild(node Node) {
	a.Children = append(a.Children, node)
}
