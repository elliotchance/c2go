package ast

type DefaultStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseDefaultStmt(line string) *DefaultStmt {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &DefaultStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *DefaultStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
