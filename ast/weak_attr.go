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

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *WeakAttr) Address() Address {
	return n.Addr
}
