package ast

type ReturnStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseReturnStmt(line string) *ReturnStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ReturnStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ReturnStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
