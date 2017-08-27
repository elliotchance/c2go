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

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CaseStmt) Address() Address {
	return n.Addr
}
