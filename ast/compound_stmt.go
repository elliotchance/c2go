package ast

type CompoundStmt struct {
	Address  string
	Position string
	Children []Node

	// TODO: remove this
	BelongsToSwitch bool
}

func parseCompoundStmt(line string) *CompoundStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &CompoundStmt{
		Address:         groups["address"],
		Position:        groups["position"],
		Children:        []Node{},
		BelongsToSwitch: false,
	}
}

func (n *CompoundStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
