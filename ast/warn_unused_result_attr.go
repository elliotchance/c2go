package ast

type WarnUnusedResultAttr struct {
	Address  string
	Position string
	Children []Node
}

func parseWarnUnusedResultAttr(line string) *WarnUnusedResultAttr {
	groups := groupsFromRegex(`<(?P<position>.*)>( warn_unused_result)?`, line)

	return &WarnUnusedResultAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *WarnUnusedResultAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
