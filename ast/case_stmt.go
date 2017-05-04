package ast

type CaseStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseCaseStmt(line string) *CaseStmt {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &CaseStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *CaseStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
