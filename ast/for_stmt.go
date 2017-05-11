package ast

type ForStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseForStmt(line string) *ForStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ForStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *ForStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
