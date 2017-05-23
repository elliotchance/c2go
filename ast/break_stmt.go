package ast

type BreakStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseBreakStmt(line string) *BreakStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &BreakStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *BreakStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
