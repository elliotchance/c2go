package ast

type WhileStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseWhileStmt(line string) *WhileStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &WhileStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *WhileStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
