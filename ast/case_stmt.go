package ast

type CaseStmt struct {
	Addr     Address
	Position string
	Children []Node
}

func parseCaseStmt(line string) *CaseStmt {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &CaseStmt{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CaseStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
