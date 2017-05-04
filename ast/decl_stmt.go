package ast

type DeclStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseDeclStmt(line string) *DeclStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &DeclStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *DeclStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
