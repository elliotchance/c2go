package ast

type CompoundStmt struct {
	Addr     Address
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
		Addr:            ParseAddress(groups["address"]),
		Position:        groups["position"],
		Children:        []Node{},
		BelongsToSwitch: false,
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CompoundStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
