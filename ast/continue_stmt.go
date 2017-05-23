package ast

type ContinueStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseContinueStmt(line string) *ContinueStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ContinueStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ContinueStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
