package ast

type SwitchStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseSwitchStmt(line string) *SwitchStmt {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &SwitchStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *SwitchStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
