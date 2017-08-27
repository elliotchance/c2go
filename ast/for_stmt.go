package ast

type ForStmt struct {
	Addr     Address
	Position string
	Children []Node
}

func parseForStmt(line string) *ForStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ForStmt{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ForStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
