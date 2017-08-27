package ast

type ModeAttr struct {
	Addr     Address
	Position string
	Name     string
	Children []Node
}

func parseModeAttr(line string) *ModeAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> (?P<name>.+)",
		line,
	)

	return &ModeAttr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Name:     groups["name"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ModeAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ModeAttr) Address() Address {
	return n.Addr
}
