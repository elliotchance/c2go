package ast

type DoStmt struct {
	Addr     Address
	Position string
	Children []Node
}

func parseDoStmt(line string) *DoStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &DoStmt{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *DoStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
