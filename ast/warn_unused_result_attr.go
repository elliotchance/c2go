package ast

type WarnUnusedResultAttr struct {
	Addr     Address
	Position string
	Children []Node
}

func parseWarnUnusedResultAttr(line string) *WarnUnusedResultAttr {
	groups := groupsFromRegex(`<(?P<position>.*)>( warn_unused_result)?`, line)

	return &WarnUnusedResultAttr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *WarnUnusedResultAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *WarnUnusedResultAttr) Address() Address {
	return n.Addr
}
