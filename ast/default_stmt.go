package ast

type DefaultStmt struct {
	Addr     Address
	Position string
	Children []Node
}

func parseDefaultStmt(line string) *DefaultStmt {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &DefaultStmt{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *DefaultStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
