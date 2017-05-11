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

func (n *ReturnStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
