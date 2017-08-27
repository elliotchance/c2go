package ast

type IfStmt struct {
	Addr     Address
	Position string
	Children []Node
}

func parseIfStmt(line string) *IfStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &IfStmt{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *IfStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
