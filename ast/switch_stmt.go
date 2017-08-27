package ast

type SwitchStmt struct {
	Addr     Address
	Position string
	Children []Node
}

func parseSwitchStmt(line string) *SwitchStmt {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &SwitchStmt{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *SwitchStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
