package ast

type DeclStmt struct {
	Addr     Address
	Position string
	Children []Node
}

func parseDeclStmt(line string) *DeclStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &DeclStmt{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *DeclStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
