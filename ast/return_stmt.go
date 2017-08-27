package ast

type ReturnStmt struct {
	Addr     Address
	Position string
	Children []Node
}

func parseReturnStmt(line string) *ReturnStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ReturnStmt{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ReturnStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
